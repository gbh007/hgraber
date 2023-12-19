package pgwithexternal

import (
	"app/internal/controller"
	"app/internal/fileStorage/externalfile"
	"app/internal/request"
	"app/internal/service/bookHandler"
	"app/internal/service/pageHandler"
	"app/internal/service/webServer"
	"app/internal/storage/postgresql"
	"app/pkg/worker"
	"context"
	"fmt"
)

type App struct {
	fs *externalfile.Storage

	ws *webServer.WebServer

	async *controller.Object
}

// Deprecated: сейчас только для теста.
func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) error {
	cfg := parseFlag()

	app.fs = externalfile.New(cfg.fs.Token, cfg.fs.Scheme, cfg.fs.Addr)
	db, err := postgresql.Connect(ctx, cfg.PGSource)
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
	requester := request.New()

	bh := bookHandler.New(bookHandler.Config{
		Storage:   db,
		Requester: requester,
		Monitor:   monitor,
	})
	ph := pageHandler.New(pageHandler.Config{
		Storage:   db,
		Files:     app.fs,
		Requester: requester,
		Monitor:   monitor,
	})

	app.ws = webServer.New(webServer.Config{
		Storage:       db,
		Book:          bh,
		Page:          ph,
		Files:         app.fs,
		Monitor:       monitor,
		Addr:          cfg.ws.Addr,
		Token:         cfg.ws.Token,
		StaticDirPath: cfg.ws.Static,
	})

	app.async = controller.NewObject()
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
