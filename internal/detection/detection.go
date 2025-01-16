package detection

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	config2 "github.com/ace-zhaoy/wireguard-helper/internal/detection/config"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection/contract"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection/plugins"
)

type Detection struct {
	config  config2.Config
	plugins []contract.Plugin
}

func NewDetection() *Detection {
	return &Detection{}
}

func (d *Detection) LoadConfig(conf config2.Config) (err error) {
	defer errors.Recover(func(e error) { err = e })
	pls, err := d.loadPlugins(conf)
	errors.Check(err)
	d.config = conf
	d.plugins = pls
	return
}

func (d *Detection) loadPlugins(conf config2.Config) (pls []contract.Plugin, err error) {
	defer errors.Recover(func(e error) { err = e })
	for _, p := range conf.Plugins {
		builder, ok := plugins.GetBuilder(p.Name)
		if !ok {
			errors.Check(errors.NewWithStack("plugin not found: %s", p.Name))
		}
		plugin, err1 := builder(p.Config)
		errors.Check(err1)
		pls = append(pls, plugin)
		log.Debug("load plugin: %s", p.Name)
	}
	return
}

func (d *Detection) Check() (result bool, err error) {
	defer errors.Recover(func(e error) { err = e })
	if len(d.plugins) == 0 {
		err = errors.NewWithStack("no plugin found")
		errors.Check(err)
	}
	for _, p := range d.plugins {
		res, err1 := p.Check()
		log.Debug("plugin check result: %v", res)
		errors.Check(err1)
		if !res {
			return
		}
	}
	result = true
	return
}
