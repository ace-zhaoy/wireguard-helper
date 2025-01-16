package entity

import "github.com/ace-zhaoy/wireguard-helper/internal/tunnel/config"

type Handler func(*Context)

type Context struct {
	OriginalConfig config.Config
	DirtyConfig    config.Config
	Err            error
	index          int
	handlers       []Handler
}

func NewContext(originalConfig, dirtyConfig config.Config, handlers []Handler) *Context {
	return &Context{
		OriginalConfig: originalConfig,
		DirtyConfig:    dirtyConfig,
		index:          -1,
		handlers:       handlers,
	}
}

func (c *Context) Next() {
	c.index++
	for c.index < len(c.handlers) {
		if c.Err != nil {
			return
		}
		c.handlers[c.index](c)
		c.index++
	}
}
