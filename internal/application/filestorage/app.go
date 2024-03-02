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

	cfg := parseFlag()

	logger := logger.New(cfg.Debug, cfg.Trace)

	logger.InfoContext(ctx, "Инициализация сервера")

	webtool := web.New(logger, cfg.Debug)
	storage := filesystem.New(cfg.LoadPath, cfg.ExportPath, cfg.ReadOnly, logger)
	controller := externalfile.New(storage, cfg.Addr, cfg.Token, logger, webtool)

	async := async.New(logger)
	async.RegisterRunner(ctx, controller)

	err := storage.Prepare(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	logger.InfoContext(ctx, "Система запущена")

	err = async.Serve(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	logger.InfoContext(ctx, "Процессы завершены, выход")
}
