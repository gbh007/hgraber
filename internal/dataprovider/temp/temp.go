package temp

import (
	"context"
	"sync"
	"time"
)

const (
	exportQueueCapacity = 100
	agentHandleTTL      = time.Minute * 10 // FIXME: в конфигурацию
)

type pageSimple struct {
	BookID     int
	PageNumber int
}

type Storage struct {
	exportQueue      []int
	exportQueueMutex *sync.Mutex

	lockBookHandle      map[int]time.Time
	lockBookHandleMutex *sync.RWMutex

	lockPageHandle      map[pageSimple]time.Time
	lockPageHandleMutex *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		exportQueue:      make([]int, 0, exportQueueCapacity),
		exportQueueMutex: new(sync.Mutex),

		lockBookHandle:      make(map[int]time.Time),
		lockBookHandleMutex: new(sync.RWMutex),
		lockPageHandle:      make(map[pageSimple]time.Time),
		lockPageHandleMutex: new(sync.RWMutex),
	}
}

func (s *Storage) AddExport(_ context.Context, bookID int) {
	s.exportQueueMutex.Lock()
	defer s.exportQueueMutex.Unlock()

	s.exportQueue = append(s.exportQueue, bookID)
}

func (s *Storage) ExportList(_ context.Context) []int {
	s.exportQueueMutex.Lock()
	defer s.exportQueueMutex.Unlock()

	// Оптимизация, чтобы не пересоздавать пустые массивы
	if len(s.exportQueue) == 0 {
		return nil
	}

	list := s.exportQueue

	// Сбрасываем очередь
	s.exportQueue = make([]int, 0, exportQueueCapacity)

	return list
}
