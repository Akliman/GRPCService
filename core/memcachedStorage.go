package core

import (
	"GRPCService/external/memcache"
	"GRPCService/logger"
	"context"
	"errors"
	"google.golang.org/grpc/status"
	"strconv"
	"sync"
)

var MemcachedCurrentStorage *MemcachedStorage

type MemcachedStorage struct {
	sync.RWMutex
}

func (s *MemcachedStorage) Get(ctx context.Context, key int32) (string, error) {
	s.RLock()
	defer s.RUnlock()

	result, err := memcache.Get(strconv.Itoa(int(key)))
	if err != nil {
		logger.LogrusLogger.Error("Error getting value. err: ", err)
		return "", err
	}

	if result == "" {
		return "", status.Error(3, "value is not found")
	}

	return result, nil
}

func (s *MemcachedStorage) Put(ctx context.Context, key int32, val string) error {
	s.Lock()
	defer s.Unlock()

	result, err := memcache.Get(strconv.Itoa(int(key)))
	if err != nil {
		logger.LogrusLogger.Error("Error getting value. err: ", err)
		return err
	}

	if result != "" {
		return status.Error(3, "key is already exist")
	}

	err = memcache.Set(strconv.Itoa(int(key)), val, 0)
	if err != nil {
		logger.LogrusLogger.Error("Error setting value. err: ", err)
		return err
	}

	return nil
}

func (s *MemcachedStorage) Delete(ctx context.Context, key int32) error {
	s.Lock()
	defer s.Unlock()

	result, err := memcache.Get(strconv.Itoa(int(key)))
	if err != nil {
		logger.LogrusLogger.Error("Error getting value. err: ", err)
		return err
	}

	if result == "" {
		return status.Error(3, "value is not found")
	}

	err = memcache.Delete(strconv.Itoa(int(key)))
	if err != nil {
		logger.LogrusLogger.Error("Error delete value. err: ", err)
		return err
	}

	return nil
}

func CreateMemcahed(ctx context.Context, adr string) error {
	if MemcachedCurrentStorage != nil {
		return errors.New("storage is already exist")
	}

	err := memcache.NewConnection(ctx, adr)
	if err != nil {
		return err
	}

	MemcachedCurrentStorage = &MemcachedStorage{}
	return nil
}
