package grpc_client

import (
	"biletter/internal/config"
	"biletter/pkg/logging"
)

type GrpcClient struct {
	//payment  router.
	logger *logging.Logger
}

func NewGrpcRouterClient(cfg config.GrpcClient, logger *logging.Logger) (*GrpcClient, error) {
	//routerDial, err := grpc.Dial(cfg.RouterService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	logger.Error(err)
	//	return nil, err
	//}
	//
	//docDial, err := grpc.Dial(cfg.DocService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	logger.Error(err)
	//	return nil, err
	//}
	//
	//fileDial, err := grpc.Dial(cfg.StorageService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	logger.Error(err)
	//	return nil, err
	//}
	//
	//gatewayDial, err := grpc.Dial(cfg.GatewayIntegrationService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//if err != nil {
	//	logger.Error(err)
	//	return nil, err
	//}

	client := GrpcClient{
		//router:  router.NewRouterClient(routerDial),
		//doc:     router.NewDocClient(docDial),
		//storage: router.NewFileStorageClient(fileDial),
		//gateway: router.NewGatewayIntegrationClient(gatewayDial),
		logger: logger,
	}

	return &client, nil
}
