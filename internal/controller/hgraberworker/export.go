package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/pkg/ctxtool"
	"context"
	"time"
)

const (
	exportWorkerInterval      = time.Second * 5
	exportWorkerQueueSize     = 100
	exportWorkerHandlersCount = 3
)

func (c *Controller) serveExportWorker(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "worker-export")

	w := worker.New[int](
		exportWorkerQueueSize,
		exportWorkerInterval,
		c.logger,
		func(ctx context.Context, bookID int) {
			err := c.useCases.ExportBook(ctx, bookID)
			if err != nil {
				c.logger.Error(ctx, err)
			}
		},
		c.useCases.ExportList,
	)

	c.register("export", w)

	w.Serve(ctx, exportWorkerHandlersCount)
}