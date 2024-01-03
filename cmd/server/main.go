package main

import (
	"app/internal/application/server"
	"app/internal/dataprovider/logger"
	"app/pkg/ctxtool"
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

	logger := logger.New(false, false)

	logger.Info(ctx, "Инициализация сервера")

	app := server.New()

	err := app.Init(ctx, logger)
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
