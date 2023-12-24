package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/hgraber"
	"app/pkg/ctxtool"
	"context"
	"time"
)

const (
	pageWorkerInterval      = time.Second * 15
	pageWorkerQueueSize     = 10000
	pageWorkerHandlersCount = 10
)

func (c *Controller) servePageWorker(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "worker-page")

	w := worker.New[hgraber.Page](
		pageWorkerQueueSize,
		pageWorkerInterval,
		c.logger,
		func(ctx context.Context, page hgraber.Page) {
			err := c.useCases.LoadPageWithUpdate(ctx, page)
			if err != nil {
				c.logger.Error(ctx, err)
			}
		},
		c.useCases.GetUnsuccessPages,
	)

	c.register("page", w)

	w.Serve(ctx, pageWorkerHandlersCount)
}
