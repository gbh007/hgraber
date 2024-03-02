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

	if !c.hasAgent {
		wg.Add(2)

		go func() {
			defer wg.Done()
			c.servePageWorker(ctx)
		}()

		go func() {
			defer wg.Done()
			c.serveBookWorker(ctx)
		}()
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		c.serveExportWorker(ctx)
	}()

	go func() {
		defer close(done)

		c.logger.InfoContext(ctx, "Запущен воркер")
		defer c.logger.InfoContext(ctx, "Воркер остановлен")

		wg.Wait()

	}()

	return done, nil
}
