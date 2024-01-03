package main

import (
	"app/internal/application/inmemory"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	inmemory.Serve(ctx)
}
