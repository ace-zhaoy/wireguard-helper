package tunnelmanager

import (
	"context"
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	"github.com/ace-zhaoy/go-utils/ucondition"
	"github.com/ace-zhaoy/go-utils/ujson"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/config"
	"time"
)

type TunnelInterface interface {
	LoadConfig(conf config.Config) (err error)
	Connect() (err error)
	PostConnectionWait()
	PostDisconnectionWait()
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
	defer t.Disconnect()
	tunnelLen, i := len(t.tunnels), 0
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		checked, err1 := t.detection.Check()
		errors.Check(err1)
		if checked {
			time.Sleep(t.config.DetectionInterval)
			continue
		}
		t.Disconnect()

		i %= tunnelLen
		j := ucondition.If(i == 0, tunnelLen-1, i-1)
		perTunnel, tun := t.tunnels[j], t.tunnels[i]
		i++
		perTunnel.PostDisconnectionWait()
		err = tun.Connect()
		errors.Check(err)
		tun.PostConnectionWait()
	}
}

func (t *TunnelManager) Disconnect() {
	connected, interfaceName, _ := t.isConnected()
	if !connected {
		return
	}
	log.Info("disconnect %s", interfaceName)
	_ = t.disconnect(interfaceName)
}
