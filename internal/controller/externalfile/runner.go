package externalfile

import (
	"app/pkg/ctxtool"
	"context"
	"errors"
	"net/http"
	"time"
)

func (*Controller) Name() string {
	return "external file controller"
}

func (c *Controller) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	webCtx := ctxtool.NewSystemContext(parentCtx, "exf-ct")
	server := c.makeServer(webCtx)

	go func() {
		defer close(done)

		c.logger.InfoContext(webCtx, "Запущен веб сервер")
		defer c.logger.InfoContext(webCtx, "Веб сервер остановлен")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			c.logger.ErrorContext(webCtx, err.Error())
		}
	}()

	go func() {
		<-parentCtx.Done()
		c.logger.InfoContext(webCtx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(webCtx), time.Second*10)
		defer cancel()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			c.logger.ErrorContext(webCtx, err.Error())
		}
	}()

	return done, nil
}
