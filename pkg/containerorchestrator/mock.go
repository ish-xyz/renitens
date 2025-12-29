package containerorchestrator

type MockClient struct {
	Nodes []*NodeInfo
}

func (m *MockClient) SetNodes(nodes []*NodeInfo) {
	m.Nodes = nodes
}

func (m *MockClient) GetNodes() ([]*NodeInfo, error) {
	return m.Nodes, nil
}
