package tasks

import (
	"context"
	"fmt"
	"sync"
)

type storageRam struct {
	statuses map[taskID]*TaskStatus
	sync.Mutex
}

func MakeStorageRam() *storageRam {
	return &storageRam{
		statuses: map[string]*TaskStatus{},
	}
}

func (s *storageRam) PutStatus(ctx context.Context, status TaskStatus) error {
	s.Lock()
	defer s.Unlock()
	s.statuses[status.ID] = &status
	return nil
}

func (s *storageRam) GetStatus(ctx context.Context, id string) (*TaskStatus, error) {
	s.Lock()
	defer s.Unlock()
	if status := s.statuses[id]; status != nil {
		copy := *status
		return &copy, nil
	}
	return nil, fmt.Errorf("not found")
}
