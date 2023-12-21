package inmemory

import (
	"app/internal/controller"
	"app/internal/fileStorage/filememory"
	"app/internal/request"
	"app/internal/service/bookHandler"
	"app/internal/service/pageHandler"
	"app/internal/service/webServer"
	"app/internal/storage/jdb"
	"app/pkg/worker"
	"app/system"
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
	system.Init(system.LogConfig{
		EnableFile:   !cfg.Log.DisableFileErr,
		AppendMode:   cfg.Log.EnableAppendFileErr,
		EnableStdErr: !cfg.Log.DisableStdErr,
	})

	// FIXME: не будет работать
	if cfg.Log.DebugMode {
		ctx = system.WithDebug(ctx)
	}

	// FIXME: не будет работать
	if cfg.Log.DebugFullpathMode {
		system.EnableFullpath(ctx)
	}

	app.async = controller.NewObject()
	app.fs = filememory.New()

	db := jdb.Init(ctx, nil)

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
		Addr:          fmt.Sprintf("%s:%d", cfg.WebServer.Host, cfg.WebServer.Port),
		Token:         cfg.WebServer.Token,
		StaticDirPath: cfg.WebServer.StaticDirPath,
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
