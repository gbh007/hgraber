package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain"
	"app/pkg/ctxtool"
	"context"
	"time"
)

const (
	bookWorkerInterval      = time.Second * 15
	bookWorkerQueueSize     = 10000
	bookWorkerHandlersCount = 10
)

func (c *Controller) serveBookWorker(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "worker-book")

	w := worker.New[domain.Book](
		bookWorkerQueueSize,
		bookWorkerInterval,
		c.logger,
		c.useCases.ParseWithUpdate,
		c.useCases.GetUnloadedBooks,
	)

	c.register("book", w)

	w.Serve(ctx, bookWorkerHandlersCount)
}
