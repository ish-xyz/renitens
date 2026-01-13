package containerorchestrator

type MockCOClient struct {
	nodes       []*Node
	podLocation string
}

func NewMockCOClient() *MockCOClient {
	return &MockCOClient{
		nodes:       []*Node{},
		podLocation: "",
	}
}

func (m *MockCOClient) SetNodes(nodes []*Node) {
	m.nodes = nodes
}

func (m *MockCOClient) GetNodes() ([]*Node, error) {
	return m.nodes, nil
}

func (m *MockCOClient) SetPodLocation(podLocation string) {
	m.podLocation = podLocation
}

func (m *MockCOClient) GetPodLocation(namespace string, pod string) (string, error) {
	return m.podLocation, nil
}
