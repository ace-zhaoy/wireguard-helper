package tunnel

import (
	"fmt"
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	"github.com/ace-zhaoy/go-utils/ujson"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/config"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/contract"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/entity"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/plugins"
	"os"
	"strings"
	"text/template"
	"time"
)

type Tunnel struct {
	config  config.Config
	plugins []contract.Plugin
}

func NewTunnel(conf config.Config) (*Tunnel, error) {
	t := &Tunnel{}
	err := t.LoadConfig(conf)
	return t, err
}

func (t *Tunnel) PostConnectionWait() {
	time.Sleep(t.config.PostConnectionWait)
}

func (t *Tunnel) PostDisconnectionWait() {
	time.Sleep(t.config.PostDisconnectionWait)
}

func (t *Tunnel) LoadConfig(conf config.Config) (err error) {
	defer errors.Recover(func(e error) { err = errors.Wrap(e, "param: %v", ujson.ToJson(conf)) })
	pls, err := t.loadPlugins(conf)
	errors.Check(err)
	t.plugins = pls
	t.config = conf
	return
}

func (t *Tunnel) loadPlugins(config config.Config) (pls []contract.Plugin, err error) {
	defer errors.Recover(func(e error) { err = e })
	for _, p := range config.Plugins {
		bd, ok := plugins.GetBuilder(p.Name)
		if !ok {
			errors.Check(errors.NewWithStack("plugin not found: " + p.Name))
		}
		pl, err1 := bd(p.Config)
		errors.Check(err1)
		pls = append(pls, pl)
		log.Debug("load plugin: %s", p.Name)
	}
	return
}

func (t *Tunnel) tplParse(tunnelName, tplFile string, data map[string]any) (configFile string, err error) {
	defer errors.Recover(func(e error) { err = e })
	tmpl, err := template.ParseFiles(tplFile)
	errors.Check(errors.WithStack(err))
	f, err := os.CreateTemp("", fmt.Sprintf("%s-*.conf", tunnelName))
	errors.Check(errors.WithStack(err))
	defer f.Close()
	err = tmpl.Execute(f, data)
	errors.Check(errors.WithStack(err))
	tplFile = f.Name()
	return tplFile, nil
}

func (t *Tunnel) connectHandler() entity.Handler {
	return func(ctx *entity.Context) {
		defer errors.Recover(func(e error) { ctx.Err = e })
		ctx.DirtyConfig.TplFile = strings.ReplaceAll(ctx.DirtyConfig.TplFile, "{name}", ctx.DirtyConfig.Name)
		dc := ctx.DirtyConfig
		log.Info("connect: %s %s", dc.Name, dc.Addr)
		dc.TplData["name"] = dc.Name
		dc.TplData["addr"] = dc.Addr
		log.Debug("tpl file %s, data: %s", ctx.DirtyConfig.TplFile, ujson.ToJson(dc.TplData))
		configFile, err := t.tplParse(dc.Name, dc.TplFile, dc.TplData)
		errors.Check(err)
		log.Debug("connect config file: %s", configFile)
		errors.Check(t.connect(configFile))
		if ctx.DirtyConfig.PostConnectionWait > 0 {
			time.Sleep(ctx.DirtyConfig.PostConnectionWait)
		}
	}
}

func (t *Tunnel) Connect() (err error) {
	defer errors.Recover(func(e error) { err = e })
	handlers := make([]entity.Handler, 0, len(t.plugins)+1)
	for _, p := range t.plugins {
		handlers = append(handlers, p.Handler)
	}
	handlers = append(handlers, t.connectHandler())
	dirtyConfig := t.config.Copy()
	log.Debug("tunnel config: %s", ujson.ToJson(t.config))
	log.Debug("tunnel dirty config: %s", ujson.ToJson(dirtyConfig))
	ctx := entity.NewContext(t.config, *dirtyConfig, handlers)
	ctx.Next()
	errors.Check(errors.WithStack(ctx.Err))
	return
}
