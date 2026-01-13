package containerorchestrator

type ContainerOrchestratorClient interface {
	GetNodes() ([]*Node, error)
	GetPodLocation(namespace string, pod string) (string, error)
}

type Node struct {
	Name        string
	IP          string
	Checkpoints []*CheckpointInfo
}

type CheckpointInfo struct {
	Timestamp int
	Namespace string
	Pod       string
	Container string
}

type Namespace string

type Pod struct {
	Name string
}
