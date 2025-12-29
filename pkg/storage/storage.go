package storage

type StorageClient interface {
	Dump() map[string]*CheckpointInfo
	Read(key string) (*CheckpointInfo, error)
	Write(key string, val *CheckpointInfo) error
	ForceWrite(key string, val *CheckpointInfo) error
}
type CheckpointInfo struct {
	Name         string   `json:"name"`
	Nodes        []string `json:"nodes"`
	CreationDate []string `json:"creationDate"`
}
