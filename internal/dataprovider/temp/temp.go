package temp

import (
	"context"
	"sync"
)

const exportQueueCapacity = 100

type Storage struct {
	exportQueue      []int
	exportQueueMutex *sync.Mutex
}

func New() *Storage {
	return &Storage{
		exportQueue:      make([]int, 0, exportQueueCapacity),
		exportQueueMutex: new(sync.Mutex),
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
