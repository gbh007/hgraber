package main

import (
	"app/internal/config"
	"app/internal/controller"
	"app/internal/service/fileStorage"
	"app/internal/service/titleHandler"
	"app/internal/service/webServer"
	"app/internal/storage/jdb"
	"app/internal/storage/sqlite"
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
		DBSaver interface {
			Save(ctx context.Context, path string, force bool) error
		}
		storage *stopwatch.Stopwatch
		err     error
	)

	controller := controller.NewObject()

	switch config.Base.DBType {
	case "jdb":
		storageJDB := jdb.Init(ctx, config.Base.DBFilePath)

		DBSaver = storageJDB
		storage = stopwatch.WithStopwatch(storageJDB)

		err := storageJDB.Load(ctx, config.Base.DBFilePath)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		controller.RegisterRunner(ctx, storageJDB)
	case "sqlite":
		sqliteDB, err := sqlite.Connect(ctx, config.Base.DBFilePath)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		err = sqliteDB.MigrateAll(ctx)
		if err != nil {
			system.Error(ctx, err)

			os.Exit(1)
		}

		storage = stopwatch.WithStopwatch(sqliteDB)

	default:
		system.Warning(ctx, "не поддерживаемый тип БД")

		os.Exit(1)

	}

	system.Info(ctx, "База загружена")

	titleService := titleHandler.Init(storage)
	pageService := fileStorage.Init(storage)

	err = system.SetFileStoragePath(ctx, config.Base.FileStoragePath)
	if err != nil {
		system.Error(ctx, err)

		os.Exit(1)
	}

	err = system.SetFileExportPath(ctx, config.Base.FileExportPath)
	if err != nil {
		system.Error(ctx, err)

		os.Exit(1)
	}

	if !config.Base.OnlyView {
		go parseTaskFile(ctx, titleService)

		controller.RegisterRunner(ctx, titleService)
		controller.RegisterRunner(ctx, pageService)
	}

	webServer := webServer.Init(storage, titleService, pageService, config.WebServer)
	controller.RegisterRunner(ctx, webServer)

	system.Info(ctx, "Завершение работы, ожидание завершения процессов")

	err = controller.Run(ctx)
	if err != nil {
		system.Error(ctx, err)

		os.Exit(1)
	}

	system.Info(ctx, "Процессы завершены")

	if DBSaver != nil {
		if DBSaver.Save(ctx, config.Base.DBFilePath, false) == nil {
			system.Info(ctx, "База сохранена")
		} else {
			system.Warning(ctx, "База не сохранена")
		}
	}

	system.Info(ctx, "Выход")
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
