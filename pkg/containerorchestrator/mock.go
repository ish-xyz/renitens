package containerorchestrator

type MockCOClient struct {
	Nodes []*NodeInfo
}

func NewMockCOClient() *MockCOClient {
	return &MockCOClient{
		Nodes: []*NodeInfo{},
	}
}

func (m *MockCOClient) SetNodes(nodes []*NodeInfo) {
	m.Nodes = nodes
}

func (m *MockCOClient) GetNodes() ([]*NodeInfo, error) {
	return m.Nodes, nil
}
