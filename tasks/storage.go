package tasks

import "sync"

type IStatusStorage interface {
	PutStatus(status TaskStatus)
	GetStatus(id string) *TaskStatus
}

func Storage() IStatusStorage {
	return ramStorage
}

var ramStorage = &MemStorage{
	statuses: map[string]*TaskStatus{},
}

type MemStorage struct {
	statuses map[taskID]*TaskStatus
	sync.Mutex
}

func (m *MemStorage) PutStatus(status TaskStatus) {
	m.Lock()
	defer m.Unlock()
	m.statuses[status.ID] = &status
}

func (m *MemStorage) GetStatus(id string) *TaskStatus {
	m.Lock()
	defer m.Unlock()
	if status := m.statuses[id]; status != nil {
		copy := *status
		return &copy
	}
	return nil
}
