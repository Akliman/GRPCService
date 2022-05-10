package handler

import (
	protos "GRPCService/api/grpc/protos"
	"GRPCService/core"
	"context"
	"google.golang.org/grpc/status"
)

type KvStorageServiceServer struct {
	protos.UnimplementedKvStorageServiceServer
}

func (s *KvStorageServiceServer) Get(ctx context.Context, req *protos.KeyRequest) (*protos.ValResponce, error) {
	storages := core.GetStorage()

	//Для локального храрилища
	if local, ok := storages.(*core.LocalStorage); ok {
		result, err := local.Get(ctx, req.GetKey())
		if err != nil {
			return &protos.ValResponce{Result: result}, err
		}
		return &protos.ValResponce{Result: result}, nil
	}

	//Для memcached харнилища
	if memc, ok := storages.(*core.MemcachedStorage); ok {
		result, err := memc.Get(ctx, req.GetKey())
		if err != nil {
			return &protos.ValResponce{Result: result}, err
		}
		return &protos.ValResponce{Result: result}, nil
	}

	return nil, status.Error(13, "Storage error")
}

func (s *KvStorageServiceServer) Put(ctx context.Context, req *protos.KeyValRequest) (*protos.Empty, error) {
	storages := core.GetStorage()

	//Для локального храрилища
	if local, ok := storages.(*core.LocalStorage); ok {
		err := local.Put(ctx, req.GetKey(), req.GetValue())
		if err != nil {
			return &protos.Empty{}, err
		}
		return &protos.Empty{}, nil
	}

	//Для memcached харнилища
	if memc, ok := storages.(*core.MemcachedStorage); ok {
		err := memc.Put(ctx, req.GetKey(), req.GetValue())
		if err != nil {
			return &protos.Empty{}, err
		}
		return &protos.Empty{}, nil
	}

	return &protos.Empty{}, status.Error(13, "Storage error")
}

func (s *KvStorageServiceServer) Delete(ctx context.Context, req *protos.KeyRequest) (*protos.Empty, error) {
	storages := core.GetStorage()

	//Для локального храрилища
	if local, ok := storages.(*core.LocalStorage); ok {
		err := local.Delete(ctx, req.GetKey())
		if err != nil {
			return &protos.Empty{}, err
		}
		return &protos.Empty{}, nil
	}

	//Для memcached харнилища
	if memc, ok := storages.(*core.MemcachedStorage); ok {
		err := memc.Delete(ctx, req.GetKey())
		if err != nil {
			return &protos.Empty{}, err
		}
		return &protos.Empty{}, nil
	}

	return &protos.Empty{}, nil
}
