package worker

import (
	"app/internal/domain"
	"slices"
	"sync"
)

type Monitor struct {
	workers map[string]domain.WorkerStat
	mutex   *sync.RWMutex
}

func NewMonitor() *Monitor {
	return &Monitor{
		workers: make(map[string]domain.WorkerStat),
		mutex:   new(sync.RWMutex),
	}
}

func (m *Monitor) Register(name string, worker domain.WorkerStat) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.workers[name] = worker
}

func (m *Monitor) Info() []domain.MonitorStat {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	res := make([]domain.MonitorStat, 0, len(m.workers))

	for name, worker := range m.workers {
		res = append(res, domain.MonitorStat{
			Name:         name,
			InQueueCount: worker.InQueueCount(),
			InWorkCount:  worker.InWorkCount(),
			RunnersCount: worker.RunnersCount(),
		})
	}

	slices.SortFunc(res, func(a, b domain.MonitorStat) int {
		if a.Name == b.Name {
			return 0
		}

		if a.Name < b.Name {
			return -1
		}

		return 1
	})

	return res
}
