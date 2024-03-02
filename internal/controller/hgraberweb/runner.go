package hgraberweb

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

	webCtx := ctxtool.NewSystemContext(parentCtx, "web")
	server := makeServer(webCtx, ws)

	go func() {
		defer close(done)

		ws.logger.InfoContext(webCtx, "Запущен веб сервер")
		defer ws.logger.InfoContext(webCtx, "Веб сервер остановлен")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ws.logger.ErrorContext(webCtx, err.Error())
		}

	}()

	go func() {
		<-parentCtx.Done()
		ws.logger.InfoContext(webCtx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(webCtx), time.Second*10)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			ws.logger.ErrorContext(webCtx, err.Error())
		}
	}()

	return done, nil
}
