package inmemory

import (
	"app/internal/controller/async"
	"app/internal/controller/bookHandler"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/pageHandler"
	"app/internal/dataprovider/fileStorage/filememory"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
	"fmt"
)

type App struct {
	fs *filememory.Storage

	ws *hgraberweb.WebServer

	async *async.Controller
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) error {
	cfg := parseFlag()

	logger := logger.New(cfg.Log.DebugMode)
	webtool := web.New(logger, cfg.Log.DebugMode)

	app.async = async.New(logger)
	app.fs = filememory.New()

	db := jdb.Init(ctx, logger, nil)

	monitor := worker.NewMonitor()
	loader := loader.New(logger)
	useCases := hgraber.New(db, logger, loader, app.fs)

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

	app.ws = hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       monitor,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
		Webtool:       webtool,
	})

	app.async.RegisterRunner(ctx, app.ws)
	app.async.RegisterRunner(ctx, bh)
	app.async.RegisterRunner(ctx, ph)

	return nil
}

func (app *App) Serve(ctx context.Context) error {
	err := app.async.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
