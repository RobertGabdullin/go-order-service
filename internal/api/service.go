package api

import (
	"context"
	"sync"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/service"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	pb.UnimplementedOrderServiceServer
	storageService service.StorageService
	logger         logger.Logger
	mu             sync.Mutex
}

func NewServer(st service.StorageService, log logger.Logger) *server {
	return &server{
		storageService: st,
		logger:         log,
	}
}

func (s *server) AcceptOrder(ctx context.Context, req *pb.AcceptOrderRequest) (*emptypb.Empty, error) {
	s.logEvent("AcceptOrder", req)

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
	s.logEvent("AcceptReturn", req)

	cmd := commands.SetAcceptReturn(s.storageService, int(req.User), int(req.Order))
	_, err := cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) DeliverOrder(ctx context.Context, req *pb.DeliverOrderRequest) (*emptypb.Empty, error) {
	s.logEvent("DeliverOrder", req)

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
	s.logEvent("GetOrders", req)

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
	s.logEvent("GetReturns", req)

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
	s.logEvent("ReturnOrder", req)

	cmd := commands.SetReturnOrd(s.storageService, int(req.Order))
	_, err := cmd.Execute(&s.mu)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
