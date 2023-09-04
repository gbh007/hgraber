package base

import (
	"app/system"
	"context"
	"net"
)

func NewBaseContext(ctx context.Context) func(l net.Listener) context.Context {
	return func(l net.Listener) context.Context { return system.NewUserContext(ctx) }
}
