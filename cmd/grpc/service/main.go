package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/api"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/config"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/event_broker"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/executor"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error loading config file: %v\n", err)
		return
	}

	connUrl := cfg.Database.URL
	if connUrl == "" {
		fmt.Println("Database URL is not set in the config file")
		return
	}
	postgresStorage, err := storage.NewOrderStorage(connUrl)
	if err != nil {
		fmt.Println(err)
		return
	}
	wrapperStorage, err := storage.NewWrapperStorage(connUrl)
	if err != nil {
		fmt.Println(err)
		return
	}

	orderService := service.NewPostgresService(postgresStorage, wrapperStorage)
	executor := executor.NewOrderCommandExecutor(orderService)

	lis, err := net.Listen("tcp", cfg.Grpc.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var kafkaClient *event_broker.KafkaClient
	if cfg.App.OutputMode == config.KafkaOutputMode {
		kafkaClient, err = event_broker.NewKafkaClient(cfg.Kafka.Brokers, nil)
		if err != nil {
			fmt.Printf("Error initializing Kafka: %v\n", err)
			return
		}
		defer kafkaClient.CloseProducer()
		go event_broker.StartConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	}

	logger := logger.KafkaLogger{
		OutputMode:  cfg.App.OutputMode,
		KafkaTopic:  cfg.Kafka.Topic,
		KafkaClient: kafkaClient,
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, api.NewServer(executor, logger))

	log.Printf("gRPC server listening on %s", lis.Addr().String())
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err = pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost"+cfg.Grpc.Port, opts)
	if err != nil {
		log.Fatalf("failed to register HTTP-gateway: %v", err)
	}

	log.Printf("HTTP-gateway server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
