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

	wg.Add(len(c.workerUnits))
	for _, w := range c.workerUnits {
		go func(ctx context.Context, w WorkerUnit) {
			defer wg.Done()
			ctx = ctxtool.NewSystemContext(ctx, "worker-"+w.Name())
			w.Serve(ctx)
		}(ctx, w)
	}

	go func() {
		defer close(done)

		c.logger.InfoContext(ctx, "Запущен воркер")
		defer c.logger.InfoContext(ctx, "Воркер остановлен")

		wg.Wait()

	}()

	return done, nil
}
