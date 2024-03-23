package simple

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/filesystem"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/logger"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/dataprovider/temp"
	"app/internal/usecase/hasher"
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/ctxtool"
	"context"
	"fmt"
)

func Serve(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "main")
	cfg := parseFlag()

	logger := logger.New(cfg.Log.Debug, cfg.Log.Trace)
	logger.InfoContext(ctx, "Инициализация сервера")

	webtool := web.New(logger, cfg.Log.Debug)

	async := async.New(logger)
	fileStorage := filesystem.New(cfg.Base.FileStoragePath, cfg.Base.FileExportPath, cfg.Base.OnlyView, logger)

	err := fileStorage.Prepare(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	storage := jdb.Init(ctx, logger, &cfg.Base.DBFilePath)

	if !cfg.Base.OnlyView {
		err = storage.Load(ctx, cfg.Base.DBFilePath)
		if err != nil {
			logger.ErrorContext(ctx, err.Error())

			return
		}

		async.RegisterRunner(ctx, storage)
		async.RegisterAfterStop(ctx, func() {
			if storage.Save(ctx, cfg.Base.DBFilePath, false) == nil {
				logger.InfoContext(ctx, "База сохранена")
			} else {
				logger.WarnContext(ctx, "База не сохранена")
			}
		})
	}

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(storage, logger, loader, fileStorage, tempStorage, false)
	hasherUC := hasher.New(storage, fileStorage)

	workerUnits := []hgraberworker.WorkerUnit{
		hgraberworker.NewExportWorkerUnit(useCases, logger),
		hgraberworker.NewHashWorkerUnit(hasherUC, logger),
		hgraberworker.NewBookWorkerUnit(useCases, logger),
		hgraberworker.NewPageWorkerUnit(useCases, logger),
	}

	worker := hgraberworker.New(logger, workerUnits)

	webServer := hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		OuterAddr:     fmt.Sprintf("http://%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
		Webtool:       webtool,
	})

	async.RegisterRunner(ctx, webServer)

	if !cfg.Base.OnlyView {
		async.RegisterRunner(ctx, worker)
	}

	logger.InfoContext(ctx, "Система запущена")

	err = async.Serve(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	logger.InfoContext(ctx, "Процессы завершены, выход")
}
