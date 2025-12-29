package containerorchestrator

type ContainerOrchestratorClient interface {
	GetNodes() ([]*NodeInfo, error)
}

type NodeInfo struct {
	Name        string
	IP          string
	Checkpoints []string
}
