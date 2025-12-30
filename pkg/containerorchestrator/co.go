package containerorchestrator

type ContainerOrchestratorClient interface {
	GetNodes() ([]*Node, error)
}

type Node struct {
	Name        string
	IP          string
	Checkpoints []string
}

type Namespace string

type Pod struct {
	Name string
}
