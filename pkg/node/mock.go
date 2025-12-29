package node

type MockNodeClient struct{}

type GetCheckpoints func() map[string]*CheckpointInfo

type CheckpointInfo struct {
	Name         string   `json:"name"`
	Nodes        []string `json:"nodes"`
	CreationDate []string `json:"creationDate"`
}

func NewMockNodeClient() *MockNodeClient {
	return &MockNodeClient{}
}
func (m *MockNodeClient) GetCheckpoints() map[string]*CheckpointInfo {
	return map[string]*CheckpointInfo{}
}
