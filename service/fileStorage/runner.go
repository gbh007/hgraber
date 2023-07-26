package fileStorage

import (
	"app/system"
	"context"
)

func (s *Service) Name() string {
	return "page storage"
}

func (s *Service) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ctx := system.NewSystemContext(parentCtx, "Page-loader")

		s.runFull(ctx)
	}()

	return done, nil
}
