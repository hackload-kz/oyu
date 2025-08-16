package grpc_controller

import (
	"biletter/internal/adapters/db/postgresql"
	pb "biletter/internal/pb/router"
	"biletter/pkg/logging"
)

type GrpcServer struct {
	pb.IntegratorServer
	logger          *logging.Logger
	postgresStorage postgresql.PostgresStorage
}

func NewIntegratorGrpcServer(logger *logging.Logger, postgresStorage postgresql.PostgresStorage) *GrpcServer {
	return &GrpcServer{
		logger:          logger,
		postgresStorage: postgresStorage,
	}
}
