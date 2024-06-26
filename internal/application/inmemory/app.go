package inmemory

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberagent"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/filememory"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/logger"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/dataprovider/temp"
	"app/internal/usecase/agentserver"
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
	logger.InfoContext(ctx, "Инициализация")

	hasAgent := cfg.Ag.Addr != ""

	webtool := web.New(logger, cfg.Log.Debug)

	async := async.New(logger)
	fileStorage := filememory.New()

	storage := jdb.Init(ctx, logger, nil)

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

	if hasAgent {
		agentUseCases := agentserver.New(logger, storage, tempStorage, fileStorage)
		agentServer := hgraberagent.New(agentUseCases, cfg.Ag.Addr, cfg.Ag.Token, logger, webtool)
		async.RegisterRunner(ctx, agentServer)
	}

	async.RegisterRunner(ctx, webServer)
	async.RegisterRunner(ctx, worker)

	logger.InfoContext(ctx, "Система запущена")

	err := async.Serve(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	logger.InfoContext(ctx, "Процессы завершены, выход")
}
