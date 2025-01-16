package logs

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog"
	"github.com/ace-zhaoy/glog/log"
	"github.com/ace-zhaoy/go-utils/ujson"
	"github.com/ace-zhaoy/wireguard-helper/pkg/utils"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level     string `json:"level"`
	File      string `json:"file"`
	Encoding  string `json:"encoding"`
	AddCaller bool   `json:"add_caller"`
}

func Init(config Config) (err error) {
	defer errors.Recover(func(e error) { err = errors.Wrap(e, "param: %s", ujson.ToJson(config)) })
	utils.SetDefaultValue(&config.Level, "info")
	lvl, err := zapcore.ParseLevel(config.Level)
	errors.Check(errors.WithStack(err))

	logConfig := glog.NewDefaultConfig()
	logConfig.AddCaller = config.AddCaller
	logConfig.FormatEnabled = true
	logConfig.Level = lvl
	if config.File != "" {
		logConfig.Core.OutputPaths = []string{config.File}
	}
	if config.Encoding != "" {
		logConfig.Core.Encoding = config.Encoding
	}

	l, err := logConfig.Build(glog.AddCallerSkip(1))
	errors.Check(err)
	log.SetLogger(l)
	return
}
