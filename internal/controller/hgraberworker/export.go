package hgraberworker

import (
	"app/internal/controller/internal/worker"
	"app/pkg/ctxtool"
	"context"
	"time"
)

const (
	exportWorkerInterval      = time.Second * 5
	exportWorkerQueueSize     = 1000
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
				c.logger.ErrorContext(ctx, err.Error())
			}
		},
		c.useCases.ExportList,
	)

	c.register("export", w)

	w.Serve(ctx, exportWorkerHandlersCount)
}
