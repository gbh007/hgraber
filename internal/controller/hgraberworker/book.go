package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/hgraber"
	"app/pkg/ctxtool"
	"context"
	"time"
)

func (c *Controller) serveBookWorker(ctx context.Context) {
	const (
		interval      = time.Second * 15
		queueSize     = 10000
		handlersCount = 10
	)

	ctx = ctxtool.NewSystemContext(ctx, "worker-book")

	w := worker.New[hgraber.Book](
		queueSize,
		interval,
		c.logger,
		c.hgraberUseCases.ParseWithUpdate,
		c.hgraberUseCases.GetUnloadedBooks,
	)

	c.register("book", w)

	w.Serve(ctx, handlersCount)
}
