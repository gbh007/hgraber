package inmemory

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/filememory"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/logger"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/dataprovider/temp"
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

	if cfg.Log.DebugMode {
		logger.SetDebug(cfg.Log.DebugMode)
	}

	webtool := web.New(logger, cfg.Log.DebugMode)

	app.async = async.New(logger)
	fileStorage := filememory.New()

	storage := jdb.Init(ctx, logger, nil)

	loader := loader.New(logger)
	tempStorage := temp.New()
	useCases := hgraber.New(storage, logger, loader, fileStorage, tempStorage, false)

	worker := hgraberworker.New(useCases, logger, false)

	webServer := hgraberweb.New(hgraberweb.Config{
		UseCases:      useCases,
		Monitor:       worker,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
		Webtool:       webtool,
	})

	app.async.RegisterRunner(ctx, webServer)
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
