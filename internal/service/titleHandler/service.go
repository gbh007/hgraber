package titleHandler

import (
	"app/internal/domain"
	"context"
	"sync"
)

type storage interface {
	GetUnloadedTitles(ctx context.Context) []domain.Title
	NewTitle(ctx context.Context, name string, URL string, loaded bool) (int, error)
	UpdateTitleAuthors(ctx context.Context, id int, authors []string) error
	UpdateTitleCategories(ctx context.Context, id int, categories []string) error
	UpdateTitleCharacters(ctx context.Context, id int, characters []string) error
	UpdateTitleGroups(ctx context.Context, id int, groups []string) error
	UpdateTitleLanguages(ctx context.Context, id int, languages []string) error
	UpdateTitleName(ctx context.Context, id int, name string) error
	UpdateTitlePages(ctx context.Context, id int, pages []domain.Page) error
	UpdateTitleParodies(ctx context.Context, id int, parodies []string) error
	UpdateTitleTags(ctx context.Context, id int, tags []string) error
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
