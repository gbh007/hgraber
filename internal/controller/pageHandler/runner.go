package pageHandler

import (
	"app/internal/domain"
	"app/pkg/ctxtool"
	"context"
)

func (s *Service) Name() string {
	return "page handler"
}

func (s *Service) Start(parentCtx context.Context) (chan struct{}, error) {
	done := make(chan struct{})

	go func() {
		defer close(done)

		ctx := ctxtool.NewSystemContext(parentCtx, "Page-loader")

		s.worker.Serve(ctx, handlersCount)
	}()

	return done, nil
}

func (s *Service) handle(ctx context.Context, page domain.Page) {
	err := s.useCases.LoadPageWithUpdate(ctx, page)
	if err != nil {
		s.logger.Error(ctx, err)
	}
}
