package server

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberagent"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/externalfile"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/logger"
	"app/internal/dataprovider/storage/postgresql"
	"app/internal/dataprovider/temp"
	"app/internal/usecase/agentserver"
	"app/internal/usecase/hasher"
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/ctxtool"
	"context"
)

func Serve(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "main")
	cfg := parseFlag()

	logger := logger.New(cfg.log.Debug, cfg.log.Trace)
	logger.InfoContext(ctx, "Инициализация сервера")

	hasAgent := cfg.ag.Addr != ""

	webtool := web.New(logger, cfg.log.Debug)
	async := async.New(logger)

	fileStorage := externalfile.New(cfg.fs.Token, cfg.fs.Scheme, cfg.fs.Addr, logger)
	storage, err := postgresql.Connect(ctx, cfg.PGSource, logger)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	if !cfg.ReadOnly {
		err = storage.MigrateAll(ctx)
		if err != nil {
			logger.ErrorContext(ctx, err.Error())

			return
		}
	}

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(storage, logger, loader, fileStorage, tempStorage, hasAgent)
	hasherUC := hasher.New(storage, fileStorage)

	workerUnits := []hgraberworker.WorkerUnit{
		hgraberworker.NewExportWorkerUnit(useCases, logger),
		hgraberworker.NewHashWorkerUnit(hasherUC, logger),
	}

	if !hasAgent {
		workerUnits = append(
			workerUnits,
			hgraberworker.NewBookWorkerUnit(useCases, logger),
			hgraberworker.NewPageWorkerUnit(useCases, logger),
		)
	}

	worker := hgraberworker.New(logger, workerUnits)

	if hasAgent && !cfg.ReadOnly {
		agentUseCases := agentserver.New(logger, storage, tempStorage, fileStorage)
		agentServer := hgraberagent.New(agentUseCases, cfg.ag.Addr, cfg.ag.Token, logger, webtool)
		async.RegisterRunner(ctx, agentServer)
	}

	webServer := hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          cfg.ws.Addr,
		OuterAddr:     cfg.ws.OuterAddr,
		Token:         cfg.ws.Token,
		StaticDirPath: cfg.ws.Static,
		Logger:        logger,
		Webtool:       webtool,
	})

	async.RegisterRunner(ctx, webServer)

	if !cfg.ReadOnly {
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
