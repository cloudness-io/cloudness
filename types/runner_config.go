package types

type RunnerConfig struct {
	EnableRunner    bool
	Engine          string
	Hostname        string
	ParallelWorkers int64
	PollingInterval int64
	MaxCPU          float64
	MaxMemory       float64
	UpdatedAt       int64

	//kubernetes
	KubeNameSpace string
}
