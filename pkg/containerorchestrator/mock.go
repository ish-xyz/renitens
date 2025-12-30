package containerorchestrator

type MockCOClient struct {
	Nodes      []*Node
	Namespaces []*Namespace
	Pods       []*Pod
}

func NewMockCOClient() *MockCOClient {
	return &MockCOClient{
		Nodes:      []*Node{},
		Namespaces: []*Namespace{},
		Pods:       []*Pod{},
	}
}

func (m *MockCOClient) SetNodes(nodes []*Node) {
	m.Nodes = nodes
}

func (m *MockCOClient) GetNodes() ([]*Node, error) {
	return m.Nodes, nil
}

func (m *MockCOClient) GetNamespaces() ([]*Namespace, error) {
	return m.Namespaces, nil
}

func (m *MockCOClient) GetPods(namespace string) ([]*Pod, error) {
	return m.Pods, nil
}
