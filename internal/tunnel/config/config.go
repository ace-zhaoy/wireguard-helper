package config

import (
	"github.com/jinzhu/copier"
	"time"
)

type Config struct {
	Name                  string         `json:"name" yaml:"name"`
	PostConnectionWait    time.Duration  `json:"post_connection_wait" yaml:"post_connection_wait"`
	PostDisconnectionWait time.Duration  `json:"post_disconnection_wait" yaml:"post_disconnection_wait"`
	Plugins               PluginConfigs  `json:"plugins" yaml:"plugins"`
	Addr                  string         `json:"addr" yaml:"addr"`
	TplFile               string         `json:"tpl_file" yaml:"tpl_file"`
	TplData               map[string]any `json:"tpl_data" yaml:"tpl_data"`
}

type PluginConfig struct {
	Name   string         `json:"name" yaml:"name"`
	Config map[string]any `json:"config" yaml:"config"`
}

type PluginConfigs []PluginConfig

func (pcs *PluginConfigs) Merge(in PluginConfig) {
	for i, existing := range *pcs {
		if existing.Name == in.Name {
			(*pcs)[i].Config = deepMerge(existing.Config, in.Config)
			return
		}
	}

	*pcs = append(*pcs, in)
}

func deepMerge(dst, src map[string]any) map[string]any {
	if dst == nil {
		dst = make(map[string]any)
	}
	for k, v := range src {
		if vMap, ok := v.(map[string]any); ok {
			if existingMap, ok := dst[k].(map[string]any); ok {
				dst[k] = deepMerge(existingMap, vMap)
			} else {
				dst[k] = deepMerge(nil, vMap)
			}
		} else {
			dst[k] = v
		}
	}
	return dst
}

func (c *Config) Copy() *Config {
	conf := new(Config)
	_ = copier.CopyWithOption(conf, c, copier.Option{DeepCopy: true})
	return conf
}

func (c *Config) Merge(conf Config) {
	if conf.Name != "" {
		c.Name = conf.Name
	}
	if conf.PostConnectionWait != 0 {
		c.PostConnectionWait = conf.PostConnectionWait
	}
	if conf.PostDisconnectionWait != 0 {
		c.PostDisconnectionWait = conf.PostDisconnectionWait
	}

	for _, plugin := range conf.Plugins {
		c.Plugins.Merge(plugin)
	}

	if conf.Addr != "" {
		c.Addr = conf.Addr
	}
	if conf.TplFile != "" {
		c.TplFile = conf.TplFile
	}
	for k, v := range conf.TplData {
		c.TplData[k] = v
	}
}
