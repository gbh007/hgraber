package main

import (
	"app/internal/config"
	"app/internal/controller"
	"app/internal/fileStorage/filememory"
	"app/internal/request"
	"app/internal/service/bookHandler"
	"app/internal/service/pageHandler"
	"app/internal/service/webServer"
	"app/internal/storage/jdb"
	"app/internal/storage/stopwatch"
	"app/pkg/worker"
	"app/system"
	"context"
	"os/signal"
	"syscall"
)

func main() {
	config := config.ParseFlag()

	system.Init(system.LogConfig{
		EnableFile:   !config.Log.DisableFileErr,
		AppendMode:   config.Log.EnableAppendFileErr,
		EnableStdErr: !config.Log.DisableStdErr,
	})

	notifyCtx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	ctx := system.NewSystemContext(notifyCtx, "Main")

	if config.Log.DebugMode {
		ctx = system.WithDebug(ctx)
	}

	if config.Log.DebugFullpathMode {
		system.EnableFullpath(ctx)
	}

	system.Debug(ctx, "Версия", system.Version)
	system.Debug(ctx, "Коммит", system.Commit)
	system.Debug(ctx, "Собрано", system.BuildAt)

	system.Info(ctx, "Инициализация базы")

	var (
		storage *stopwatch.Stopwatch
		err     error
	)

	controller := controller.NewObject()

	storageJDB := jdb.Init(ctx, config.Base.DBFilePath)

	storage = stopwatch.WithStopwatch(storageJDB)

	system.Info(ctx, "База загружена")

	fStor := filememory.New()

	monitor := worker.NewMonitor()
	requester := request.New()

	titleService := bookHandler.Init(storage, requester, monitor)
	pageService := pageHandler.Init(storage, fStor, requester, monitor)

	controller.RegisterRunner(ctx, titleService)
	controller.RegisterRunner(ctx, pageService)

	webServer := webServer.Init(storage, titleService, pageService, fStor, monitor, config.WebServer)
	controller.RegisterRunner(ctx, webServer)

	system.Info(ctx, "Система запущена")

	err = controller.Run(ctx)
	if err != nil {
		system.Error(ctx, err)
	}

	system.Info(ctx, "Процессы завершены, выход")
}
