package contract

import (
	"github.com/ace-zhaoy/wireguard-helper/internal/tunnel/entity"
)

type Plugin interface {
	Handler(ctx *entity.Context)
}
