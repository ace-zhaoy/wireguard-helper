package config

import "github.com/ace-zhaoy/glog/log"

type LogHook struct {
}

func NewLogHook() *LogHook {
	return &LogHook{}
}

func (h *LogHook) Notify(configName string, err error) {
	log.Error("configName: %s, error: %+v", configName, err)
}
