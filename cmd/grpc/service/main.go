package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/config"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/storage"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	storageService service.StorageService
	mu             sync.Mutex
}

func (s *server) AcceptOrder(ctx context.Context, req *pb.AcceptOrderRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) AcceptReturn(ctx context.Context, req *pb.AcceptReturnRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) DeliverOrder(ctx context.Context, req *pb.DeliverOrderRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *server) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	return nil, nil
}

func (s *server) GetReturns(ctx context.Context, req *pb.GetReturnsRequest) (*pb.GetReturnsResponse, error) {
	return nil, nil
}

func (s *server) ReturnOrder(ctx context.Context, req *pb.ReturnOrderRequest) (*emptypb.Empty, error) {
	return nil, nil
}

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

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, &server{
		storageService: orderService,
	})

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
