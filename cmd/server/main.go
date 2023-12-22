package main

import (
	"app/internal/application/server"
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

	logger.Info(ctx, "Инициализация сервера")

	app := server.New()

	err := app.Init(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	logger.Info(ctx, "Система запущена")

	err = app.Serve(ctx)
	if err != nil {
		logger.Error(ctx, err)
	}

	logger.Info(ctx, "Процессы завершены, выход")
}
