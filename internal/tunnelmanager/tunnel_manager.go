package tunnelmanager

import (
	"context"
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	"github.com/ace-zhaoy/go-utils/ujson"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/config"
	"time"
)

type TunnelInterface interface {
	LoadConfig(ctx context.Context, conf config.Config) (err error)
	Connect(ctx context.Context) (err error)
	Disconnect(ctx context.Context) (err error)
	PostConnectionWait(ctx context.Context)
	PostDisconnectionWait(ctx context.Context)
}

type TunnelManager struct {
	config    Config
	tunnels   []TunnelInterface
	detection *detection.Detection
}

func NewTunnelManager(conf Config, dc *detection.Detection) (tm *TunnelManager, err error) {
	t := &TunnelManager{}
	t.LoadDetection(dc)
	err = t.LoadConfig(conf)
	if err != nil {
		return
	}
	tm = t
	return
}

func (t *TunnelManager) LoadDetection(detection *detection.Detection) {
	t.detection = detection
}

func (t *TunnelManager) LoadConfig(conf Config) (err error) {
	defer errors.Recover(func(e error) { err = errors.Wrap(e, "param: %v", ujson.ToJson(conf)) })
	tunnels := make([]TunnelInterface, 0, len(conf.Tunnels))
	for _, v := range conf.Tunnels {
		addrs := v.Addrs
		if len(addrs) == 0 {
			addrs = []string{v.Addr}
		}
		for _, v2 := range addrs {
			v.Addr = v2
			confCopy := conf.DefaultConfig.Copy()
			confCopy.Merge(v.Config)
			log.Debug("load tunnel: %s %s", confCopy.Name, confCopy.Addr)
			tun, err1 := tunnel.NewTunnel(*confCopy)
			errors.Check(err1)
			tunnels = append(tunnels, tun)
		}
	}
	if len(tunnels) == 0 {
		errors.Check(errors.New("no tunnel"))
	}
	t.config, t.tunnels = conf, tunnels
	return
}

func (t *TunnelManager) Connect(ctx context.Context) (err error) {
	defer errors.Recover(func(e error) { err = e })

	t.Disconnect()
	defer t.Disconnect()

	tunnelLen, i := len(t.tunnels), 0
	checkChan := make(chan bool, 0)
	defer close(checkChan)

	checkFunc := func() bool {
		checked, err1 := t.detection.Check()
		if err1 != nil {
			log.Error("detection check error: %+v", err1)
			checked = true
		}
		return checked
	}

	var tun TunnelInterface

	for {
		go func() {
			checked := checkFunc()
			checkChan <- checked
		}()

		select {
		case <-ctx.Done():
			if tun != nil {
				_ = tun.Disconnect(ctx)
			}
			return
		case checked := <-checkChan:
			if checked {
				select {
				case <-ctx.Done():
					if tun != nil {
						_ = tun.Disconnect(ctx)
					}
					return
				case <-time.After(t.config.DetectionInterval):
				}
				continue
			}
		}

		if tun != nil {
			_ = tun.Disconnect(ctx)
			tun.PostDisconnectionWait(ctx)
		}

		i %= tunnelLen
		err = t.tunnels[i].Connect(ctx)
		errors.Check(err)
		t.tunnels[i].PostConnectionWait(ctx)

		tun = t.tunnels[i]
		i++
	}
}

func (t *TunnelManager) Disconnect() {
	connected, interfaceName, _ := t.isConnected()
	if !connected {
		return
	}

	log.Info("disconnect %s", interfaceName)

	err := t.disconnect(interfaceName)
	if err != nil {
		log.Error("disconnect error: %v", err)
	}
}
