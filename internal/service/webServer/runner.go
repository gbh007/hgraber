package webServer

import (
	"app/pkg/ctxtool"
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

	webCtx := ctxtool.NewSystemContext(parentCtx, "Web-srv")
	server := makeServer(webCtx, ws)

	go func() {
		defer close(done)

		ws.logger.Info(webCtx, "Запущен веб сервер")
		defer ws.logger.Info(webCtx, "Веб сервер остановлен")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ws.logger.Error(webCtx, err)
		}

	}()

	go func() {
		<-parentCtx.Done()
		ws.logger.Info(webCtx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(webCtx), time.Second*10)
		defer cancel()

		ws.logger.IfErr(webCtx, server.Shutdown(shutdownCtx))
	}()

	return done, nil
}
