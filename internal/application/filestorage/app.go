package filestorage

import (
	"app/internal/controller"
	"app/internal/controller/externalfile"
	"app/internal/dataprovider/fileStorage/filesystem"
	"app/internal/usecase/web"
	"app/pkg/logger"
	"context"
	"fmt"
)

type App struct {
	storage    *filesystem.Storage
	controller *externalfile.Controller
	async      *controller.Object
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) {
	cfg := parseFlag()

	debug := false // FIXME: управлять отладкой с конфигурации

	logger := logger.New(debug)
	webtool := web.New(logger, debug)

	app.storage = filesystem.New(cfg.LoadPath, cfg.ExportPath, cfg.ReadOnly)

	app.controller = externalfile.New(app.storage, cfg.Addr, cfg.Token, logger, webtool)

	app.async = controller.NewObject(logger)
	app.async.RegisterRunner(ctx, app.controller)
}

func (app *App) Serve(ctx context.Context) error {
	err := app.storage.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	err = app.async.Run(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
