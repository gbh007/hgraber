package hgraberworker

import (
	"app/internal/domain/hgraber"
	"slices"
)

func (c *Controller) register(name string, worker hgraber.WorkerStat) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.workers[name] = worker
}

func (c *Controller) Info() []hgraber.MonitorStat {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	res := make([]hgraber.MonitorStat, 0, len(c.workers))

	for name, worker := range c.workers {
		res = append(res, hgraber.MonitorStat{
			Name:         name,
			InQueueCount: worker.InQueueCount(),
			InWorkCount:  worker.InWorkCount(),
			RunnersCount: worker.RunnersCount(),
		})
	}

	slices.SortFunc(res, func(a, b hgraber.MonitorStat) int {
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
