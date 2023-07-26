package titleHandler

import (
	"app/storage/schema"
	"app/super"
	"sync"
)

type Service struct {
	Storage super.Storage

	titleQueue              chan schema.Title
	inWorkRunnersCount      int
	inWorkRunnersCountMutex *sync.RWMutex

	asyncPathWG *sync.WaitGroup
}

func Init(storage super.Storage) *Service {
	return &Service{
		Storage:                 storage,
		titleQueue:              make(chan schema.Title, titleQueueSize),
		inWorkRunnersCountMutex: &sync.RWMutex{},
		asyncPathWG:             &sync.WaitGroup{},
	}
}
