package inmemory

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/filememory"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/dataprovider/temp"
	"app/internal/usecase/hgraber"
	"app/internal/usecase/web"
	"app/pkg/logger"
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

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(db, logger, loader, app.fs, tempStorage)

	worker := hgraberworker.New(useCases, logger)

	app.ws = hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
		Webtool:       webtool,
	})

	app.async.RegisterRunner(ctx, app.ws)
	app.async.RegisterRunner(ctx, worker)

	return nil
}

func (app *App) Serve(ctx context.Context) error {
	err := app.async.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
