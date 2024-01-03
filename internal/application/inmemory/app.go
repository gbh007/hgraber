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
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/ctxtool"
	"context"
	"fmt"
)

func Serve(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "main")
	logger := logger.New(false, false)
	logger.Info(ctx, "Инициализация")

	cfg := parseFlag()

	if cfg.Log.DebugMode {
		logger.SetDebug(cfg.Log.DebugMode)
	}

	hasAgent := cfg.Ag.Addr != ""

	webtool := web.New(logger, cfg.Log.DebugMode)

	async := async.New(logger)
	fileStorage := filememory.New()

	storage := jdb.Init(ctx, logger, nil)

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(storage, logger, loader, fileStorage, tempStorage, hasAgent)

	worker := hgraberworker.New(useCases, logger, hasAgent)

	webServer := hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
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

	logger.Info(ctx, "Система запущена")

	err := async.Serve(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	logger.Info(ctx, "Процессы завершены, выход")
}
