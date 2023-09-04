package fileStorage

import (
	"app/internal/domain"
	"context"
	"sync"
)

type storage interface {
	GetTitle(ctx context.Context, id int) (domain.Title, error)
	GetUnsuccessedPages(ctx context.Context) []domain.PageFullInfo
	UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error
}

type Service struct {
	storage     storage
	queue       chan qPage
	inWork      int
	inWorkMutex *sync.RWMutex

	asyncPathWG *sync.WaitGroup
}

func Init(storage storage) *Service {
	return &Service{
		storage:     storage,
		queue:       make(chan qPage, pageQueueSize),
		inWorkMutex: &sync.RWMutex{},
		asyncPathWG: &sync.WaitGroup{},
	}
}
