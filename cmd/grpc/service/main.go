package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
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
	expire, err := time.Parse("2006-01-02T15", req.Expire)
	if err != nil {
		return nil, err
	}

	cmd := commands.SetAcceptOrd(s.storageService, int(req.User), int(req.Order), int(req.Weight), int(req.BasePrice), expire, req.Wrapper)
	_, err = cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) AcceptReturn(ctx context.Context, req *pb.AcceptReturnRequest) (*emptypb.Empty, error) {
	cmd := commands.SetAcceptReturn(s.storageService, int(req.User), int(req.Order))
	_, err := cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeliverOrder(ctx context.Context, req *pb.DeliverOrderRequest) (*emptypb.Empty, error) {
	orders := make([]int, len(req.Orders))
	for i, order := range req.Orders {
		orders[i] = int(order)
	}

	cmd := commands.SetDeliverOrd(s.storageService, orders)
	_, err := cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	cmd := commands.SetGetOrds(s.storageService, int(req.User), int(req.Count))

	orders, err := cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	orderList := []*pb.Order{}
	for _, order := range orders {
		orderList = append(orderList, &pb.Order{
			Id:        int32(order.Id),
			Recipient: int32(order.Recipient),
			Expire:    order.Expire.Format("2006-01-02T15"),
			Status:    order.Status,
		})
	}

	return &pb.GetOrdersResponse{Orders: orderList}, nil
}

func (s *server) GetReturns(ctx context.Context, req *pb.GetReturnsRequest) (*pb.GetReturnsResponse, error) {
	returns, err := s.storageService.GetReturns(int(req.Offset), int(req.Limit))
	if err != nil {
		return nil, err
	}

	returnList := []*pb.Return{}
	for _, ret := range returns {
		returnList = append(returnList, &pb.Return{
			Id:         int32(ret.Id),
			Recipient:  int32(ret.Recipient),
			Expire:     ret.Expire.Format("2006-01-02T15"),
			ReturnedAt: ret.ReturnedAt.Format("2006-01-02T15"),
		})
	}

	return &pb.GetReturnsResponse{Returns: returnList}, nil
}

func (s *server) ReturnOrder(ctx context.Context, req *pb.ReturnOrderRequest) (*emptypb.Empty, error) {
	cmd := commands.SetReturnOrd(s.storageService, int(req.Order))
	_, err := cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
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
	err = pb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("failed to register HTTP-gateway: %v", err)
	}

	log.Printf("HTTP-gateway server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
