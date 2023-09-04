package main

import (
	"app/internal/config"
	"app/internal/controller"
	"app/internal/service/fileStorage"
	"app/internal/service/titleHandler"
	"app/internal/service/webServer"
	"app/internal/storage/jdb"
	"app/internal/storage/stopwatch"
	"app/system"
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
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

	mainContext := system.NewSystemContext(notifyCtx, "Main")

	if config.Log.DebugMode {
		mainContext = system.WithDebug(mainContext)
	}

	if config.Log.DebugFullpathMode {
		system.EnableFullpath(mainContext)
	}

	system.Debug(mainContext, "Версия", system.Version)
	system.Debug(mainContext, "Коммит", system.Commit)
	system.Debug(mainContext, "Собрано", system.BuildAt)

	system.Info(mainContext, "Инициализация базы")

	storageJDB := jdb.Init(mainContext, config.Base.DBFilePath)
	storage := stopwatch.WithStopwatch(storageJDB)

	err := storageJDB.Load(mainContext, config.Base.DBFilePath)
	if err != nil {
		os.Exit(1)
	}

	system.Info(mainContext, "База загружена")

	titleService := titleHandler.Init(storage)
	pageService := fileStorage.Init(storage)

	controller := controller.NewObject()
	controller.RegisterRunner(mainContext, storageJDB)

	err = system.SetFileStoragePath(mainContext, config.Base.FileStoragePath)
	if err != nil {
		os.Exit(2)
	}

	err = system.SetFileExportPath(mainContext, config.Base.FileExportPath)
	if err != nil {
		os.Exit(3)
	}

	if !config.Base.OnlyView {
		go parseTaskFile(mainContext, titleService)

		controller.RegisterRunner(mainContext, titleService)
		controller.RegisterRunner(mainContext, pageService)
	}

	webServer := webServer.Init(storage, titleService, pageService, config.WebServer)
	controller.RegisterRunner(mainContext, webServer)

	system.Info(mainContext, "Завершение работы, ожидание завершения процессов")

	err = controller.Run(mainContext)
	if err != nil {
		os.Exit(4)
	}

	system.Info(mainContext, "Процессы завершены")

	if storageJDB.Save(mainContext, config.Base.DBFilePath, false) == nil {
		system.Info(mainContext, "База сохранена")
	} else {
		system.Warning(mainContext, "База не сохранена")
	}

	system.Info(mainContext, "Выход")
}

func parseTaskFile(ctx context.Context, service *titleHandler.Service) {
	f, err := os.Open("task.txt")
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			system.Error(ctx, err)
		}
		return
	}
	defer system.IfErrFunc(ctx, f.Close)

	var (
		data []string
	)

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		u := sc.Text()
		if u != "" {
			data = append(data, u)
		}
	}

	res := service.FirstHandleMultiple(ctx, data)

	system.Info(ctx,
		fmt.Sprintf(
			"всего: %d загружено: %d дубликаты: %d ошибки: %d",
			res.TotalCount, res.LoadedCount, res.DuplicateCount, res.ErrorCount,
		))
}
