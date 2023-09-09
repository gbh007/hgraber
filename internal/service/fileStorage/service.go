package fileStorage

import (
	"app/internal/domain"
	"app/pkg/worker"
	"context"
)

type storage interface {
	GetTitle(ctx context.Context, id int) (domain.Title, error)
	GetUnsuccessedPages(ctx context.Context) []domain.PageFullInfo
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
}

type Service struct {
	storage storage

	worker *worker.Worker[qPage]
}

func Init(storage storage) *Service {
	s := &Service{
		storage: storage,
	}

	s.worker = worker.New[qPage](
		queueSize,
		interval,
		s.handle,
		s.getter,
	)

	return s
}
