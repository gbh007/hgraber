package titleHandler

import (
	"app/internal/domain"
	"app/pkg/worker"
	"context"
)

type storage interface {
	GetUnloadedTitles(ctx context.Context) []domain.Title
	NewTitle(ctx context.Context, name string, URL string, loaded bool) (int, error)
	UpdateTitlePages(ctx context.Context, id int, pages []domain.Page) error

	UpdateTitleName(ctx context.Context, id int, name string) error
	UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error
}

type Service struct {
	storage storage

	worker *worker.Worker[domain.Title]
}

func Init(storage storage) *Service {
	s := &Service{
		storage: storage,
	}

	s.worker = worker.New[domain.Title](
		titleQueueSize,
		titleInterval,
		s.updateForWorker,
		s.storage.GetUnloadedTitles,
	)

	return s
}
