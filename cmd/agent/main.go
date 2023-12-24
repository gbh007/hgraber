package main

import (
	"app/internal/application/agent"
	"app/pkg/ctxtool"
	"app/pkg/logger"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	notifyCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	ctx := ctxtool.NewSystemContext(notifyCtx, "main")

	// FIXME: сейчас 2 логгера
	logger := logger.New(false)

	app := agent.New()

	app.Init(ctx)

	logger.Info(ctx, "Система запущена")

	err := app.Serve(ctx)
	if err != nil {
		logger.Error(ctx, err)
	}

	logger.Info(ctx, "Процессы завершены, выход")
}
