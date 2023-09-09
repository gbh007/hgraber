package webServer

import (
	"app/system"
	"context"
	"errors"
	"net/http"
	"time"
)

func (ws *WebServer) Name() string {
	return "web server"
}

func (ws *WebServer) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	webCtx := system.NewSystemContext(parentCtx, "Web-srv")
	server := makeServer(webCtx, ws)

	go func() {
		defer close(done)

		system.Info(webCtx, "Запущен веб сервер")
		defer system.Info(webCtx, "Веб сервер остановлен")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			system.Error(webCtx, err)
		}

	}()

	go func() {
		<-parentCtx.Done()
		system.Info(webCtx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(webCtx), time.Second*10)
		defer cancel()

		system.IfErr(webCtx, server.Shutdown(shutdownCtx))
	}()

	return done, nil
}
