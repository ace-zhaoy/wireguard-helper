package ipincrement

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/contract"
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/entity"
	"github.com/ace-zhaoy/wireguard-helper/pkg/utils"
	"net"
	"time"
)

type Config struct {
	StartTime      time.Time     `json:"start_time"`
	IncrementCycle time.Duration `json:"increment_cycle"`
}

type IPIncrement struct {
	config Config
}

func Builder(conf map[string]any) (contract.Plugin, error) {
	return NewIPIncrement(conf)
}

func NewIPIncrement(conf map[string]any) (i *IPIncrement, err error) {
	defer errors.Recover(func(e error) { err = e })
	var config Config
	err = utils.ToStruct(conf, &config)
	errors.Check(errors.WithStack(err))
	i = &IPIncrement{config: config}
	return
}

func (i *IPIncrement) Handler(ctx *entity.Context) {
	parsedIP := net.ParseIP(ctx.DirtyConfig.Addr)
	if parsedIP == nil {
		return
	}
	ipBytes := parsedIP.To4()
	if ipBytes == nil {
		return
	}
	step := time.Now().Sub(i.config.StartTime) / i.config.IncrementCycle
	// may overflow
	ipBytes[3] += byte(step)
	ctx.DirtyConfig.Addr = ipBytes.String()
}
