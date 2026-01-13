package storage

import "fmt"

type MemStorage struct {
	// key should be -> namespace/pod/container
	stateByKey  map[string]*CheckpointInfo
	stateByNode map[string]*CheckpointInfo
}

func NewMemStorage() *MemStorage {
	{
		return &MemStorage{
			stateByKey:  make(map[string]*CheckpointInfo),
			stateByNode: make(map[string]*CheckpointInfo),
		}
	}
}

func (m *MemStorage) Dump() map[string]*CheckpointInfo {
	return m.stateByKey
}

func (m *MemStorage) Read(key string) (*CheckpointInfo, error) {
	if val, ok := m.stateByKey[key]; ok {
		return val, nil
	}
	return nil, fmt.Errorf("failed to retrieve value, key does not exist.")
}

func (m *MemStorage) Write(key string, val *CheckpointInfo) error {
	return m.write(key, val, false)
}

func (m *MemStorage) ForceWrite(key string, val *CheckpointInfo) error {
	return m.write(key, val, true)
}

func (m *MemStorage) write(key string, val *CheckpointInfo, overwrite bool) error {
	if _, ok := m.stateByKey[key]; ok {
		if overwrite {
			m.stateByKey[key] = val
			return nil
		}
		return fmt.Errorf("key already exists")
	}

	m.stateByKey[key] = val

	return nil
}
