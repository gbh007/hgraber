package server

import (
	"app/internal/controller"
	"app/internal/controller/bookHandler"
	"app/internal/controller/pageHandler"
	"app/internal/controller/webServer"
	"app/internal/dataprovider/fileStorage/externalfile"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/storage/postgresql"
	"app/internal/usecase"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
	"fmt"
)

type App struct {
	fs *externalfile.Storage

	ws *webServer.WebServer

	async *controller.Object
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) error {
	cfg := parseFlag()

	logger := logger.New(false) //FIXME

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

	monitor := worker.NewMonitor()
	loader := loader.New(logger)
	useCases := usecase.New(db, logger, loader, app.fs)

	bh := bookHandler.New(bookHandler.Config{
		UseCases: useCases,
		Monitor:  monitor,
		Logger:   logger,
	})
	ph := pageHandler.New(pageHandler.Config{
		UseCases: useCases,
		Monitor:  monitor,
		Logger:   logger,
	})

	app.ws = webServer.New(webServer.Config{
		UseCases:      useCases,
		Monitor:       monitor,
		Addr:          cfg.ws.Addr,
		Token:         cfg.ws.Token,
		StaticDirPath: cfg.ws.Static,
		Logger:        logger,
	})

	app.async = controller.NewObject(logger)
	app.async.RegisterRunner(ctx, app.ws)

	if !cfg.ReadOnly {
		app.async.RegisterRunner(ctx, bh)
		app.async.RegisterRunner(ctx, ph)
	}

	return nil
}

func (app *App) Serve(ctx context.Context) error {
	err := app.async.Run(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
