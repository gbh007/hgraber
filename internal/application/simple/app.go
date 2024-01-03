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
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/ctxtool"
	"context"
	"fmt"
)

func Serve(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "main")
	logger := logger.New(false, false)
	logger.Info(ctx, "Инициализация сервера")

	cfg := parseFlag()

	if cfg.Log.DebugMode {
		logger.SetDebug(cfg.Log.DebugMode)
	}

	webtool := web.New(logger, cfg.Log.DebugMode)

	async := async.New(logger)
	fileStorage := filesystem.New(cfg.Base.FileStoragePath, cfg.Base.FileExportPath, cfg.Base.OnlyView, logger)

	err := fileStorage.Prepare(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	storage := jdb.Init(ctx, logger, &cfg.Base.DBFilePath)

	if !cfg.Base.OnlyView {
		err = storage.Load(ctx, cfg.Base.DBFilePath)
		if err != nil {
			logger.Error(ctx, err)

			return
		}

		async.RegisterRunner(ctx, storage)
		async.RegisterAfterStop(ctx, func() {
			if storage.Save(ctx, cfg.Base.DBFilePath, false) == nil {
				logger.Info(ctx, "База сохранена")
			} else {
				logger.Warning(ctx, "База не сохранена")
			}
		})
	}

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(storage, logger, loader, fileStorage, tempStorage, false)

	worker := hgraberworker.New(useCases, logger, false)

	webServer := hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
		Webtool:       webtool,
	})

	async.RegisterRunner(ctx, webServer)

	if !cfg.Base.OnlyView {
		async.RegisterRunner(ctx, worker)
	}

	logger.Info(ctx, "Система запущена")

	err = async.Serve(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	logger.Info(ctx, "Процессы завершены, выход")
}
