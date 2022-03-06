package web

import (
	"app/system"
	"app/web/base"
	"context"
	"errors"
	"net/http"
	"time"
)

// Run запускает веб сервер
func Run(parentCtx context.Context, addr string) {
	ctx := system.NewSystemContext(parentCtx, "WEB-SRV")
	mux := http.NewServeMux()

	// обработчик статики
	mux.Handle("/", http.FileServer(http.Dir("./static")))
	mux.Handle("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir(system.GetFileStoragePath(ctx)))))

	// API
	base.AddHandler(mux, "/info", MainInfo())
	base.AddHandler(mux, "/new", NewTitle())
	base.AddHandler(mux, "/title/list", TitleList())
	base.AddHandler(mux, "/title/details", TitleInfo())
	base.AddHandler(mux, "/title/page", TitlePage())
	base.AddHandler(mux, "/to-zip", SaveToZIP())

	server := http.Server{
		Addr:        addr,
		Handler:     mux,
		ErrorLog:    system.StdErrorLogger(ctx),
		BaseContext: base.NewBaseContext(ctx),
	}

	system.AddWaiting(ctx)
	go func() {
		defer system.DoneWaiting(ctx)
		system.Info(ctx, "Запущен веб сервер")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			system.Error(ctx, err)
		}

		system.Info(ctx, "Веб сервер остановлен")
	}()

	go func() {
		<-ctx.Done()
		system.Info(ctx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		system.IfErr(ctx, server.Shutdown(shutdownCtx))
	}()
}
