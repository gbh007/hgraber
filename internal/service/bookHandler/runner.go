package bookHandler

import (
	"app/system"
	"context"
	"time"
)

const (
	titleInterval      = time.Second * 15
	titleQueueSize     = 10000
	titleHandlersCount = 10
)

func (s *Service) Name() string {
	return "title handler"
}

func (s *Service) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ctx := system.NewSystemContext(parentCtx, "Title-loader")

		s.worker.Serve(ctx, titleHandlersCount)
	}()

	return done, nil
}
