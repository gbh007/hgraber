package webServer

import (
	"app/service/webServer/base"
	"app/service/webServer/static"
	"app/super"
	"app/system"
	"context"
	"errors"
	"net/http"
	"time"
)

type WebServer struct {
	Storage   super.Storage
	Title     super.TitleHandler
	Page      super.PageHandler
	Addr      string
	StaticDir string
	Token     string
}

// Start запускает веб сервер
func Start(parentCtx context.Context, ws *WebServer) {
	ctx := system.NewSystemContext(parentCtx, "Web-srv")
	mux := http.NewServeMux()

	// обработчик статики
	if ws.StaticDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(ws.StaticDir)))
	} else {
		mux.Handle("/", http.FileServer(http.FS(static.StaticDir)))
	}

	// обработчик файлов
	mux.Handle("/file/", base.TokenHandler(ws.Token,
		http.StripPrefix(
			"/file/",
			http.FileServer(http.Dir(system.GetFileStoragePath(ctx))),
		),
	))

	// API
	mux.Handle("/auth/login", ws.routeLogin(ws.Token))
	mux.Handle("/info", base.TokenHandler(ws.Token, ws.routeMainInfo()))
	mux.Handle("/new", base.TokenHandler(ws.Token, ws.routeNewTitle()))
	mux.Handle("/title/list", base.TokenHandler(ws.Token, ws.routeTitleList()))
	mux.Handle("/title/details", base.TokenHandler(ws.Token, ws.routeTitleInfo()))
	mux.Handle("/title/page", base.TokenHandler(ws.Token, ws.routeTitlePage()))
	mux.Handle("/to-zip", base.TokenHandler(ws.Token, ws.routeSaveToZIP()))
	mux.Handle("/app/info", base.TokenHandler(ws.Token, ws.routeAppInfo()))
	mux.Handle("/title/rate", base.TokenHandler(ws.Token, ws.routeSetTitleRate()))
	mux.Handle("/title/page/rate", base.TokenHandler(ws.Token, ws.routeSetPageRate()))

	server := http.Server{
		Addr:        ws.Addr,
		Handler:     base.PanicDefender(mux),
		ErrorLog:    system.StdErrorLogger(ctx),
		BaseContext: base.NewBaseContext(ctx),
	}

	system.AddWaiting(ctx)
	go func() {
		defer system.DoneWaiting(ctx)
		system.Info(ctx, "Запущен веб сервер")
		defer system.Info(ctx, "Веб сервер остановлен")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			system.Error(ctx, err)
		}

	}()

	go func() {
		<-ctx.Done()
		system.Info(ctx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		system.IfErr(ctx, server.Shutdown(shutdownCtx))
	}()
}
