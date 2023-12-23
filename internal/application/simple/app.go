package simple

import (
	"app/internal/controller"
	"app/internal/controller/bookHandler"
	"app/internal/controller/pageHandler"
	"app/internal/controller/webServer"
	"app/internal/dataprovider/fileStorage/filesystem"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/storage/jdb"
	"app/internal/usecase/hgraber"
	"app/pkg/logger"
	"app/pkg/worker"
	"context"
	"fmt"
)

type App struct {
	fs *filesystem.Storage

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

	app.ws = webServer.New(webServer.Config{
		UseCases:      useCases,
		Monitor:       monitor,
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
		Logger:        logger,
	})

	app.async.RegisterRunner(ctx, app.ws)

	if !cfg.Base.OnlyView {
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
