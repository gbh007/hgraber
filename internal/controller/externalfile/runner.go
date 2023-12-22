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

		c.logger.Info(webCtx, "Запущен веб сервер")
		defer c.logger.Info(webCtx, "Веб сервер остановлен")

		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			c.logger.Error(webCtx, err)
		}

	}()

	go func() {
		<-parentCtx.Done()
		c.logger.Info(webCtx, "Остановка веб сервера")

		shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(webCtx), time.Second*10)
		defer cancel()

		c.logger.IfErr(webCtx, server.Shutdown(shutdownCtx))
	}()

	return done, nil
}
