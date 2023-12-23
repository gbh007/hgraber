package simple

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/filesystem"
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
	fs *filesystem.Storage

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
	app.fs = filesystem.New(cfg.Base.FileStoragePath, cfg.Base.FileExportPath, cfg.Base.OnlyView)

	err := app.fs.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	db := jdb.Init(ctx, logger, &cfg.Base.DBFilePath)

	if !cfg.Base.OnlyView {
		err = db.Load(ctx, cfg.Base.DBFilePath)
		if err != nil {
			return fmt.Errorf("app: %w", err)
		}

		app.async.RegisterRunner(ctx, db)
		app.async.RegisterAfterStop(ctx, func() {
			if db.Save(ctx, cfg.Base.DBFilePath, false) == nil {
				logger.Info(ctx, "База сохранена")
			} else {
				logger.Warning(ctx, "База не сохранена")
			}
		})
	}

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

	if !cfg.Base.OnlyView {
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
