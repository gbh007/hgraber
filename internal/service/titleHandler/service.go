package titleHandler

import (
	"app/internal/domain"
	"app/pkg/worker"
	"context"
)

type storage interface {
	GetUnloadedBooks(ctx context.Context) []domain.Book
	NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error)
	UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error

	UpdateBookName(ctx context.Context, id int, name string) error
	UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error
}

type Service struct {
	storage storage

	worker *worker.Worker[domain.Book]
}

func Init(storage storage) *Service {
	s := &Service{
		storage: storage,
	}

	s.worker = worker.New[domain.Book](
		titleQueueSize,
		titleInterval,
		s.updateForWorker,
		s.storage.GetUnloadedBooks,
	)

	return s
}
