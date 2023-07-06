package webServer

import (
	"app/service/webServer/base"
	"app/service/webServer/static"
	"app/system"
	"context"
	"errors"
	"net/http"
	"time"
)

// Start запускает веб сервер
func Start(parentCtx context.Context, addr string, staticDir string, token string) {
	ctx := system.NewSystemContext(parentCtx, "Web-srv")
	mux := http.NewServeMux()

	// обработчик статики
	if staticDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(staticDir)))
	} else {
		mux.Handle("/", http.FileServer(http.FS(static.StaticDir)))
	}

	// обработчик файлов
	mux.Handle("/file/", base.TokenHandler(token,
		http.StripPrefix(
			"/file/",
			http.FileServer(http.Dir(system.GetFileStoragePath(ctx))),
		),
	))

	// API
	mux.Handle("/auth/login", Login(token))
	mux.Handle("/info", base.TokenHandler(token, MainInfo()))
	mux.Handle("/new", base.TokenHandler(token, NewTitle()))
	mux.Handle("/title/list", base.TokenHandler(token, TitleList()))
	mux.Handle("/title/details", base.TokenHandler(token, TitleInfo()))
	mux.Handle("/title/page", base.TokenHandler(token, TitlePage()))
	mux.Handle("/to-zip", base.TokenHandler(token, SaveToZIP()))
	mux.Handle("/app/info", base.TokenHandler(token, AppInfo()))
	mux.Handle("/title/rate", base.TokenHandler(token, SetTitleRate()))
	mux.Handle("/title/page/rate", base.TokenHandler(token, SetPageRate()))

	server := http.Server{
		Addr:        addr,
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
