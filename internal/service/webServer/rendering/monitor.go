package rendering

import "app/internal/domain"

type Monitor struct {
	Workers []WorkerUnit `json:"workers"`
}

type WorkerUnit struct {
	Name    string `json:"name"`
	InQueue int    `json:"in_queue"`
	InWork  int    `json:"in_work"`
	Runners int    `json:"runners"`
}

func MonitorFromDomain(workers []domain.MonitorStat) Monitor {
	workersOut := make([]WorkerUnit, len(workers))

	convertSlice(workersOut, workers, WorkerUnitFromDomain)

	return Monitor{
		Workers: workersOut,
	}
}

func WorkerUnitFromDomain(worker domain.MonitorStat) WorkerUnit {
	return WorkerUnit{
		Name:    worker.Name,
		InQueue: worker.InQueueCount,
		InWork:  worker.InWorkCount,
		Runners: worker.RunnersCount,
	}
}
