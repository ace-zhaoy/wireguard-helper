package tunnelmanager

import (
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/config"
	"time"
)

type Config struct {
	Mode               string         `json:"mode" yaml:"mode"`
	DetectionInterval  time.Duration  `json:"detection_interval" yaml:"detection_interval"`
	ConnectTunnelNames []string       `json:"connect_tunnel_names" yaml:"connect_tunnel_names"`
	DefaultConfig      config.Config  `json:"default_config" yaml:"default_config"`
	Tunnels            []TunnelConfig `json:"tunnels" yaml:"tunnels"`
}

type TunnelConfig struct {
	config.Config `json:",squash" yaml:",inline"`
	Addrs         []string `json:"addrs" yaml:"addrs"`
}
