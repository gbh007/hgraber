package simple

import (
	"app/internal/controller/async"
	"app/internal/controller/hgraberweb"
	"app/internal/controller/hgraberworker"
	"app/internal/dataprovider/fileStorage/filesystem"
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
	fileStorage := filesystem.New(cfg.Base.FileStoragePath, cfg.Base.FileExportPath, cfg.Base.OnlyView, logger)

	err := fileStorage.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	storage := jdb.Init(ctx, logger, &cfg.Base.DBFilePath)

	if !cfg.Base.OnlyView {
		err = storage.Load(ctx, cfg.Base.DBFilePath)
		if err != nil {
			return fmt.Errorf("app: %w", err)
		}

		app.async.RegisterRunner(ctx, storage)
		app.async.RegisterAfterStop(ctx, func() {
			if storage.Save(ctx, cfg.Base.DBFilePath, false) == nil {
				logger.Info(ctx, "База сохранена")
			} else {
				logger.Warning(ctx, "База не сохранена")
			}
		})
	}

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
