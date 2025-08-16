package main

import (
	_ "biletter/docs"
	"biletter/internal/adapters/db/postgresql"
	"biletter/internal/config"
	"biletter/internal/controller/grpc_controller"
	"biletter/internal/controller/rest_controller"
	"biletter/internal/grpc_client"
	pb "biletter/internal/pb/router"
	"biletter/internal/services"
	"biletter/pkg/client/postgresql_client"
	"biletter/pkg/cors"
	"biletter/pkg/logging"
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	"github.com/julienschmidt/httprouter"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

// @title Biletter API
// @version 1.0
// @description API Server Biletter

// @host 0.0.0.0:8081
// @BasePath /

func main() {
	logger := logging.GetLogger()
	logger.Info("create router...")

	cfg := config.GetConfig()
	cfgPostgres := cfg.Storage
	cfgGrpc := cfg.GrpcClient

	router := httprouter.New()

	redisCache, err := redisConnect(cfg)
	if err != nil {
		logger.Fatalf("Redis connection was refused: %v", err)
	}

	routerClient, err := grpc_client.NewGrpcRouterClient(cfgGrpc, logger)
	if err != nil {
		logger.Fatal(err)
	}

	postgresClient, err := postgresql_client.NewClient(context.TODO(), 3, cfgPostgres)
	if err != nil {
		logger.Fatal(err)
	}

	postgresStorage := postgresql.NewPostgresStorage(postgresClient, logger, *routerClient)

	grpcServer := grpc_controller.NewIntegratorGrpcServer(logger, postgresStorage)

	IntegratorService := services.NewService(postgresStorage, logger, redisCache, *routerClient)

	restHandler := rest_controller.NewRouterHandler(postgresStorage, logger, *routerClient, IntegratorService)
	restHandler.Register(router)

	go runGrpc(grpcServer, cfg)
	runRest(router, cfg)
}

func runGrpc(grpcServer *grpc_controller.GrpcServer, cfg *config.Config) {
	logger := logging.GetLogger()
	list, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Grpc.BindIP, cfg.Grpc.Port))
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	logger.Infof("GRPC server is listening %s:%s", cfg.Grpc.BindIP, cfg.Grpc.Port)

	pb.RegisterIntegratorServer(s, grpcServer)
	if err = s.Serve(list); err != nil {
		logger.Fatal(err)
	}
}

func runRest(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}

		logger.Info("create socket...")
		socketPath := path.Join(appDir, "app.sock")
		logger.Debugf("socket created, path - %s", socketPath)

		logger.Infof("listnet unex socket...")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)

	} else {
		logger.Infof("listnet ctp...")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		panic(listenErr)
	}

	corsSettings := cors.GetCorsSettings(cfg)

	header := corsSettings.Handler(router)

	server := &http.Server{
		Handler:      header,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}

func redisConnect(cfg *config.Config) (*cache.Cache, error) {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server": fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port),
		},
	})

	_, err := ring.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	myCache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})

	return myCache, nil
}
