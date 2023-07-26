package fileStorage

import (
	"app/super"
	"sync"
)

type Service struct {
	Storage     super.Storage
	queue       chan qPage
	inWork      int
	inWorkMutex *sync.RWMutex

	asyncPathWG *sync.WaitGroup
}

func Init(storage super.Storage) *Service {
	return &Service{
		Storage:     storage,
		queue:       make(chan qPage, pageQueueSize),
		inWorkMutex: &sync.RWMutex{},
		asyncPathWG: &sync.WaitGroup{},
	}
}
