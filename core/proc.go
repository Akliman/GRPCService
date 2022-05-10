package core

import (
	"context"
)

type KVStorage interface {
	Get(ctx context.Context, key int32) (string, error)
	Put(ctx context.Context, key int32, val string) error
	Delete(ctx context.Context, key int32) error
}

func GetStorage() interface{} {
	if LocalCurrentStorage != nil {
		return LocalCurrentStorage
	}

	if MemcachedCurrentStorage != nil {
		return MemcachedCurrentStorage
	}

	return nil
}
