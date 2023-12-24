package agent

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/agent"
	"app/pkg/ctxtool"
	"context"
	"time"
)

const (
	bookWorkerInterval      = time.Second * 20
	bookWorkerQueueSize     = 60
	bookWorkerHandlersCount = 3
)

func (c *Controller) serveBookWorker(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "book")

	w := worker.New[agent.BookToHandle](
		bookWorkerQueueSize,
		bookWorkerInterval,
		c.logger,
		c.useCases.BookHandle,
		c.useCases.Books,
	)

	w.Serve(ctx, bookWorkerHandlersCount)
}
