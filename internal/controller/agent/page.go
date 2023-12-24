package agent

import (
	"app/internal/controller/internal/worker"
	"app/internal/domain/agent"
	"app/pkg/ctxtool"
	"context"
	"time"
)

const (
	pageWorkerInterval      = time.Second * 20
	pageWorkerQueueSize     = 100
	pageWorkerHandlersCount = 10
)

func (c *Controller) servePageWorker(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "page")

	w := worker.New[agent.PageToHandle](
		pageWorkerQueueSize,
		pageWorkerInterval,
		c.logger,
		c.useCases.PageHandle,
		c.useCases.Pages,
	)

	w.Serve(ctx, pageWorkerHandlersCount)
}
