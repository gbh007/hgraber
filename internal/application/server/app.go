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
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"context"
	"fmt"
)

type App struct {
	async *async.Controller
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context, logger *logger.Logger) error {
	cfg := parseFlag()

	debug := false // FIXME: получать из конфигурации
	hasAgent := cfg.ag.Addr != ""

	if debug {
		logger.SetDebug(debug)
	}

	webtool := web.New(logger, debug)
	app.async = async.New(logger)

	fileStorage := externalfile.New(cfg.fs.Token, cfg.fs.Scheme, cfg.fs.Addr, logger)
	storage, err := postgresql.Connect(ctx, cfg.PGSource, logger)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	if !cfg.ReadOnly {
		err = storage.MigrateAll(ctx)
		if err != nil {
			return fmt.Errorf("app: %w", err)
		}
	}

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(storage, logger, loader, fileStorage, tempStorage, hasAgent)

	worker := hgraberworker.New(useCases, logger, hasAgent)

	if hasAgent && !cfg.ReadOnly {
		agentUseCases := agentserver.New(logger, storage, tempStorage, fileStorage)
		agentServer := hgraberagent.New(agentUseCases, cfg.ag.Addr, cfg.ag.Token, logger, webtool)
		app.async.RegisterRunner(ctx, agentServer)
	}

	webServer := hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          cfg.ws.Addr,
		Token:         cfg.ws.Token,
		StaticDirPath: cfg.ws.Static,
		Logger:        logger,
		Webtool:       webtool,
	})

	app.async.RegisterRunner(ctx, webServer)

	if !cfg.ReadOnly {
		app.async.RegisterRunner(ctx, worker)
	}

	return nil
}

func (app *App) Serve(ctx context.Context) error {
	err := app.async.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
