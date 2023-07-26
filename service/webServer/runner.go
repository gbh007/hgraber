package webServer

import (
	"app/system"
	"context"
)

func (ws *WebServer) Name() string {
	return "web server"
}

func (ws *WebServer) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ctx := system.NewSystemContext(parentCtx, "Web server")

		// FIXME
		Start(ctx, ws)
	}()

	return done, nil
}
