package core

import (
	"context"
	"errors"
	"google.golang.org/grpc/status"
	"sync"
)

var LocalCurrentStorage *LocalStorage

type LocalStorage struct {
	sync.RWMutex
	data map[int32]string
}

func (s *LocalStorage) Get(ctx context.Context, key int32) (string, error) {
	s.RLock()
	defer s.RUnlock()

	if val, isFound := s.data[key]; isFound {
		return val, nil
	}

	return "", status.Error(3, "value is not found")
}

func (s *LocalStorage) Put(ctx context.Context, key int32, val string) error {
	s.Lock()
	defer s.Unlock()

	if _, isFound := s.data[key]; isFound {
		return status.Error(3, "key is already exist")
	}

	s.data[key] = val

	return nil
}

func (s *LocalStorage) Delete(ctx context.Context, key int32) error {
	s.Lock()
	defer s.Unlock()

	if _, isFound := s.data[key]; !isFound {
		return status.Error(3, "value is not found")
	}

	delete(s.data, key)
	return nil
}

func CreateLocal() error {
	if LocalCurrentStorage != nil {
		return errors.New("storage is already exist")
	}

	LocalCurrentStorage = &LocalStorage{data: make(map[int32]string)}
	return nil
}
