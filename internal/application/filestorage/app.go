package filestorage

import (
	"app/internal/controller"
	"app/internal/controller/externalfile"
	"app/internal/fileStorage/filesystem"
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

	app.storage = filesystem.New(cfg.LoadPath, cfg.ExportPath, cfg.ReadOnly)

	app.controller = externalfile.New(app.storage, cfg.Addr, cfg.Token)

	app.async = controller.NewObject()
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
