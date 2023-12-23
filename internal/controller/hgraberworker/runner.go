package hgraberworker

import (
	"app/pkg/ctxtool"
	"context"
	"sync"
)

func (c *Controller) Name() string {
	return "worker"
}

func (c *Controller) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	ctx := ctxtool.NewSystemContext(parentCtx, "worker")

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.servePageWorker(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.serveBookWorker(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.serveExportWorker(ctx)
	}()

	go func() {
		defer close(done)

		c.logger.Info(ctx, "Запущен воркер")
		defer c.logger.Info(ctx, "Воркер остановлен")

		wg.Wait()

	}()

	return done, nil
}
