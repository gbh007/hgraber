package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/hgraber"
	"app/pkg/ctxtool"
	"context"
	"time"
)

func (c *Controller) servePageHasher(ctx context.Context) {
	const (
		interval      = time.Second * 15
		queueSize     = 10000
		handlersCount = 10
	)

	ctx = ctxtool.NewSystemContext(ctx, "worker-hasher")

	w := worker.New[hgraber.Page](
		queueSize,
		interval,
		c.logger,
		func(ctx context.Context, page hgraber.Page) {
			err := c.hasherUseCases.HandlePage(ctx, page)
			if err != nil {
				c.logger.ErrorContext(ctx, err.Error())
			}
		},
		c.hasherUseCases.UnHashedPages,
	)

	c.register("hasher", w)

	w.Serve(ctx, handlersCount)
}
