package main

import (
	"app/internal/application/server"
	"app/system"
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

	ctx := system.NewSystemContext(notifyCtx, "main")

	system.Info(ctx, "Инициализация сервера")

	app := server.New()

	err := app.Init(ctx)
	if err != nil {
		system.Error(ctx, err)

		return
	}

	system.Info(ctx, "Система запущена")

	err = app.Serve(ctx)
	if err != nil {
		system.Error(ctx, err)
	}

	system.Info(ctx, "Процессы завершены, выход")
}
