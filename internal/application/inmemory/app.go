package inmemory

import (
	"app/internal/controller"
	"app/internal/fileStorage/filememory"
	"app/internal/request"
	"app/internal/service/bookHandler"
	"app/internal/service/pageHandler"
	"app/internal/service/webServer"
	"app/internal/storage/jdb"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
	"fmt"
)

type App struct {
	fs *filememory.Storage

	ws *webServer.WebServer

	async *controller.Object
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) error {
	cfg := parseFlag()

	logger := logger.New(cfg.Log.DebugMode)

	app.async = controller.NewObject(logger)
	app.fs = filememory.New()

	db := jdb.Init(ctx, logger, nil)

	monitor := worker.NewMonitor()
	requester := request.New(logger)

	bh := bookHandler.New(bookHandler.Config{
		Storage:   db,
		Requester: requester,
		Monitor:   monitor,
		Logger:    logger,
	})
	ph := pageHandler.New(pageHandler.Config{
		Storage:   db,
		Files:     app.fs,
		Requester: requester,
		Monitor:   monitor,
		Logger:    logger,
	})

	app.ws = webServer.New(webServer.Config{
		Storage:       db,
		Book:          bh,
		Page:          ph,
		Files:         app.fs,
		Monitor:       monitor,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
	})

	app.async.RegisterRunner(ctx, app.ws)
	app.async.RegisterRunner(ctx, bh)
	app.async.RegisterRunner(ctx, ph)

	return nil
}

func (app *App) Serve(ctx context.Context) error {
	err := app.async.Run(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
