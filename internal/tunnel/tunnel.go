package tunnel

import (
	"context"
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
	name    string
	config  config.Config
	plugins []contract.Plugin
}

func NewTunnel(conf config.Config) (*Tunnel, error) {
	t := &Tunnel{}
	err := t.LoadConfig(context.Background(), conf)
	return t, err
}

func (t *Tunnel) PostConnectionWait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-time.After(t.config.PostConnectionWait):
	}
}

func (t *Tunnel) PostDisconnectionWait(ctx context.Context) {
	select {
	case <-ctx.Done():
	case <-time.After(t.config.PostDisconnectionWait):
	}
}

func (t *Tunnel) LoadConfig(_ context.Context, conf config.Config) (err error) {
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
	dir, err := os.MkdirTemp("", "wireguard-helper-tunnel")
	errors.CheckWithStack(err)
	configFile = fmt.Sprintf("%s/%s.conf", dir, tunnelName)
	f, err := os.Create(configFile)
	errors.CheckWithStack(err)
	defer f.Close()
	err = tmpl.Execute(f, data)
	errors.Check(errors.WithStack(err))
	return
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
		t.name = dc.Name
		configFile, err := t.tplParse(dc.Name, dc.TplFile, dc.TplData)
		errors.Check(err)
		log.Debug("connect config file: %s", configFile)
		errors.Check(t.connect(configFile))
		if ctx.DirtyConfig.PostConnectionWait > 0 {
			time.Sleep(ctx.DirtyConfig.PostConnectionWait)
		}
	}
}

func (t *Tunnel) Connect(ctx context.Context) (err error) {
	defer errors.Recover(func(e error) { err = e })
	handlers := make([]entity.Handler, 0, len(t.plugins)+1)
	for _, p := range t.plugins {
		handlers = append(handlers, p.Handler)
	}
	handlers = append(handlers, t.connectHandler())
	waitChan := make(chan error)
	defer close(waitChan)

	go func() {
		dirtyConfig := t.config.Copy()
		log.Debug("tunnel config: %s", ujson.ToJson(t.config))
		log.Debug("tunnel dirty config: %s", ujson.ToJson(dirtyConfig))
		entCtx := entity.NewContext(t.config, *dirtyConfig, handlers)
		entCtx.Next()
		waitChan <- errors.WithStack(entCtx.Err)

		select {
		case <-ctx.Done():
			err1 := t.disconnect(dirtyConfig.Name)
			if err1 != nil {
				log.Error("disconnect error: %v", err1)
			}
		default:
		}
	}()

	select {
	case <-ctx.Done():
	case err1 := <-waitChan:
		errors.Check(err1)
	}
	return
}

func (t *Tunnel) Disconnect(_ context.Context) (err error) {
	log.Info("disconnect: %s", t.name)
	return t.disconnect(t.name)
}
