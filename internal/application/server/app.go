package server

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/externalfile"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/storage/postgresql"
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/logger"
	"context"
	"fmt"
)

type App struct {
	fs *externalfile.Storage

	ws *hgraberweb.WebServer

	async *async.Controller
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) error {
	cfg := parseFlag()

	debug := false // FIXME: получать из конфигурации

	logger := logger.New(debug)
	webtool := web.New(logger, debug)

	app.fs = externalfile.New(cfg.fs.Token, cfg.fs.Scheme, cfg.fs.Addr, logger)
	db, err := postgresql.Connect(ctx, cfg.PGSource, logger)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	if !cfg.ReadOnly {
		err = db.MigrateAll(ctx)
		if err != nil {
			return fmt.Errorf("app: %w", err)
		}
	}

	loader := loader.New(logger)
	useCases := hgraber.New(db, logger, loader, app.fs)

	worker := hgraberworker.New(useCases, logger)

	app.ws = hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          cfg.ws.Addr,
		Token:         cfg.ws.Token,
		StaticDirPath: cfg.ws.Static,
		Logger:        logger,
		Webtool:       webtool,
	})

	app.async = async.New(logger)
	app.async.RegisterRunner(ctx, app.ws)

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
