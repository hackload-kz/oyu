package services

import (
	"biletter/internal/adapters/db/postgresql"
	"biletter/internal/domain/entity"
	"biletter/internal/grpc_client"
	"biletter/pkg/logging"
	"context"
	"github.com/go-redis/cache/v9"
)

type Service interface {
	GetEventList(ctx context.Context, query entity.EventListQuery) (response []entity.EventForList, err error)
}

type service struct {
	logger          *logging.Logger
	postgresStorage postgresql.PostgresStorage
	redisCache      *cache.Cache
	rpc             grpc_client.GrpcClient
}

func NewService(postgresStorage postgresql.PostgresStorage, logger *logging.Logger, redisCache *cache.Cache,
	router grpc_client.GrpcClient,
) Service {
	return &service{
		logger:          logger,
		postgresStorage: postgresStorage,
		redisCache:      redisCache,
		rpc:             router,
	}
}
