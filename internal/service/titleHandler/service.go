package titleHandler

import (
	"app/internal/domain"
	"context"
	"sync"
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

	titleQueue              chan domain.Title
	inWorkRunnersCount      int
	inWorkRunnersCountMutex *sync.RWMutex

	asyncPathWG *sync.WaitGroup
}

func Init(storage storage) *Service {
	return &Service{
		storage:                 storage,
		titleQueue:              make(chan domain.Title, titleQueueSize),
		inWorkRunnersCountMutex: &sync.RWMutex{},
		asyncPathWG:             &sync.WaitGroup{},
	}
}
