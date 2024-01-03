package filestorage

import (
	"app/internal/controller/async"
	"app/internal/controller/externalfile"
	"app/internal/dataprovider/fileStorage/filesystem"
	"app/internal/dataprovider/logger"
	"app/internal/usecase/web"
	"app/pkg/ctxtool"
	"context"
)

func Serve(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "main")

	logger := logger.New(false, false)

	logger.Info(ctx, "Инициализация сервера")
	cfg := parseFlag()

	debug := false // FIXME: управлять отладкой с конфигурации

	if debug {
		logger.SetDebug(debug)
	}

	webtool := web.New(logger, debug)
	storage := filesystem.New(cfg.LoadPath, cfg.ExportPath, cfg.ReadOnly, logger)
	controller := externalfile.New(storage, cfg.Addr, cfg.Token, logger, webtool)

	async := async.New(logger)
	async.RegisterRunner(ctx, controller)

	err := storage.Prepare(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	logger.Info(ctx, "Система запущена")

	err = async.Serve(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	logger.Info(ctx, "Процессы завершены, выход")
}
