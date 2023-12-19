package main

import (
	"app/internal/application/filestorage"
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

	app := filestorage.New()
	app.Init(ctx)

	system.Info(ctx, "Система запущена")

	err := app.Serve(ctx)
	if err != nil {
		system.Error(ctx, err)
	}

	system.Info(ctx, "Процессы завершены, выход")
}
