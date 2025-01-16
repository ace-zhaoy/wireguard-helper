package plugins

import (
	"github.com/ace-zhaoy/wireguard-helper/internal/detection/contract"
	"github.com/ace-zhaoy/wireguard-helper/internal/detection/plugins/ping"
)

type Builder func(conf map[string]any) (contract.Plugin, error)

var _builder = map[string]Builder{}

func Register(name string, builder Builder) {
	_builder[name] = builder
}

func GetBuilder(name string) (Builder, bool) {
	builder, ok := _builder[name]
	return builder, ok
}

func init() {
	Register("ping", ping.Builder)
}
