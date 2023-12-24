package hgraber

type WorkerStat interface {
	InQueueCount() int
	InWorkCount() int
	RunnersCount() int
}

type MonitorStat struct {
	Name         string
	InQueueCount int
	InWorkCount  int
	RunnersCount int
}
