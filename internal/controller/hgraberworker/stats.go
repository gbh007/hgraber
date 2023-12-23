package hgraberworker

import (
	"app/internal/domain"
	"slices"
)

func (c *Controller) register(name string, worker domain.WorkerStat) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.workers[name] = worker
}

func (c *Controller) Info() []domain.MonitorStat {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	res := make([]domain.MonitorStat, 0, len(c.workers))

	for name, worker := range c.workers {
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
