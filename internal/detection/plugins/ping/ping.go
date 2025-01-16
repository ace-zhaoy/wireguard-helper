package ping

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	"github.com/ace-zhaoy/go-utils/ujson"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection/contract"
	"github.com/ace-zhaoy/wireguard-helper/pkg/utils"
	"github.com/go-ping/ping"
	"time"
)

type Config struct {
	Count      int           `json:"count" yaml:"count"`
	TTL        int           `json:"ttl" yaml:"ttl"`
	Interval   time.Duration `json:"interval" yaml:"interval"`
	Timeout    time.Duration `json:"timeout" yaml:"timeout"`
	Privileged *bool         `json:"privileged" yaml:"privileged"`
	PacketLoss float64       `json:"packet_loss" yaml:"packet_loss"`
	Targets    []Target      `json:"targets" yaml:"targets"`
	Retries    int           `json:"retries" yaml:"retries"`
}

type Target struct {
	Addr string `json:"addr" yaml:"addr"`
}

func Builder(conf map[string]any) (contract.Plugin, error) {
	return NewPing(conf)
}

type Ping struct {
	config *Config
}

func NewPing(conf map[string]any) (p *Ping, err error) {
	defer errors.Recover(func(e error) { err = e })
	var config Config
	err = utils.ToStruct(conf, &config)
	errors.Check(err)
	p = &Ping{config: &config}
	return
}

func (p *Ping) buildPinger(target Target) (pinger *ping.Pinger, err error) {
	defer errors.Recover(func(e error) { err = errors.Wrap(e, "param: %v", ujson.ToJson(target)) })
	pinger, err = ping.NewPinger(target.Addr)
	errors.Check(errors.WithStack(err))
	utils.SetValue(&pinger.Count, p.config.Count, 3)
	utils.SetValue(&pinger.TTL, p.config.TTL, 128)
	utils.SetValue(&pinger.Interval, p.config.Interval, time.Millisecond*100)
	utils.SetValue(&pinger.Timeout, p.config.Timeout, time.Second*2)
	if p.config.Privileged != nil {
		pinger.SetPrivileged(*p.config.Privileged)
	}
	return
}

func (p *Ping) check() (result bool, err error) {
	defer errors.Recover(func(e error) { err = e })
	if len(p.config.Targets) == 0 {
		errors.Check(errors.NewWithStack("targets is empty"))
	}
	for _, target := range p.config.Targets {
		log.Debug("ping target: %+v", target)
		pinger, err := p.buildPinger(target)
		errors.Check(errors.Wrap(err, "target: %v", ujson.ToJson(target)))
		err = pinger.Run()
		errors.Check(errors.Wrap(err, "target: %v", ujson.ToJson(target)))
		stats := pinger.Statistics()
		log.Debug("ping result: %+v", stats)
		if stats.PacketLoss >= p.config.PacketLoss {
			return false, nil
		}
	}
	return true, nil
}

func (p *Ping) Check() (result bool, err error) {
	for i := 0; i < p.config.Retries+1; i++ {
		result, err = p.check()
		if result {
			return
		}
	}
	return
}
