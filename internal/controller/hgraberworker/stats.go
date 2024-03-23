package hgraberworker

import (
	"app/internal/domain/hgraber"
	"slices"
)

func (c *Controller) Info() []hgraber.MonitorStat {
	res := make([]hgraber.MonitorStat, 0, len(c.workerUnits))

	for _, worker := range c.workerUnits {
		res = append(res, hgraber.MonitorStat{
			Name:         worker.Name(),
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
