package config

type Plugin struct {
	Name   string         `json:"name" yaml:"name"`
	Config map[string]any `json:"config" yaml:"config"`
}

type Config struct {
	Plugins []Plugin `json:"plugins" yaml:"plugins"`
}
