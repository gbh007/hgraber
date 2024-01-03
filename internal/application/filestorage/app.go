package filestorage

import (
	"app/internal/controller/async"
	"app/internal/controller/externalfile"
	"app/internal/dataprovider/fileStorage/filesystem"
	"app/internal/dataprovider/logger"
	"app/internal/usecase/web"
	"context"
	"fmt"
)

type App struct {
	storage    *filesystem.Storage
	controller *externalfile.Controller
	async      *async.Controller
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context, logger *logger.Logger) {
	cfg := parseFlag()

	debug := false // FIXME: управлять отладкой с конфигурации

	if debug {
		logger.SetDebug(debug)
	}

	webtool := web.New(logger, debug)

	app.storage = filesystem.New(cfg.LoadPath, cfg.ExportPath, cfg.ReadOnly, logger)

	app.controller = externalfile.New(app.storage, cfg.Addr, cfg.Token, logger, webtool)

	app.async = async.New(logger)
	app.async.RegisterRunner(ctx, app.controller)
}

func (app *App) Serve(ctx context.Context) error {
	err := app.storage.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	err = app.async.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
