package storage

import "fmt"

type MemStorage struct {
	// key should be -> namespace/pod/container
	storeForCheckpoints map[string]*CheckpointInfo
}

func NewMemStorage() *MemStorage {
	{
		return &MemStorage{
			storeForCheckpoints: make(map[string]*CheckpointInfo),
		}
	}
}

func (m *MemStorage) Dump() map[string]*CheckpointInfo {
	return m.storeForCheckpoints
}

func (m *MemStorage) Read(key string) (*CheckpointInfo, error) {
	if val, ok := m.storeForCheckpoints[key]; ok {
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
	if _, ok := m.storeForCheckpoints[key]; ok {
		if overwrite {
			m.storeForCheckpoints[key] = val
			return nil
		}
		return fmt.Errorf("key already exists")
	}

	m.storeForCheckpoints[key] = val

	return nil
}
