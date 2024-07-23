package api

import (
	"context"
	"time"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/executor"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/logger"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedOrderServiceServer
	executor executor.CommandExecutorGrpc
	logger   logger.Logger
}

func NewServer(executor executor.CommandExecutorGrpc, log logger.Logger) *Server {
	return &Server{
		executor: executor,
		logger:   log,
	}
}

func (s *Server) AcceptOrder(ctx context.Context, req *pb.AcceptOrderRequest) (*pb.AcceptOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	s.logEvent("AcceptOrder", req)

	expire, err := time.Parse("2006-01-02T15", req.Expire)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid expire format: %v", err)
	}

	cmd := commands.NewAcceptOrder(int(req.User), int(req.Order), int(req.Weight), int(req.BasePrice), expire, req.Wrapper)
	_, err = s.executor.AcceptOrder(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute command: %v", err)
	}

	return &pb.AcceptOrderResponse{}, nil
}

func (s *Server) AcceptReturn(ctx context.Context, req *pb.AcceptReturnRequest) (*pb.AcceptReturnResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	s.logEvent("AcceptReturn", req)

	cmd := commands.NewAcceptReturn(int(req.User), int(req.Order))
	_, err := s.executor.AcceptReturn(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute command: %v", err)
	}

	return &pb.AcceptReturnResponse{}, nil
}

func (s *Server) DeliverOrder(ctx context.Context, req *pb.DeliverOrderRequest) (*pb.DeliverOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	s.logEvent("DeliverOrder", req)

	orders := make([]int, len(req.Orders))
	for i, order := range req.Orders {
		orders[i] = int(order)
	}

	cmd := commands.NewDeliverOrder(orders)
	_, err := s.executor.DeliverOrder(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute command: %v", err)
	}

	return &pb.DeliverOrderResponse{}, nil
}

func (s *Server) GetOrders(ctx context.Context, req *pb.GetOrdersRequest) (*pb.GetOrdersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	s.logEvent("GetOrders", req)

	cmd := commands.NewGetOrders(int(req.User), int(req.Count))
	orders, err := s.executor.GetOrders(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
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

func (s *Server) GetReturns(ctx context.Context, req *pb.GetReturnsRequest) (*pb.GetReturnsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	s.logEvent("GetReturns", req)

	cmd := commands.NewGetReturns(int(req.Offset), int(req.Limit))
	returns, err := s.executor.GetReturns(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get returns: %v", err)
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

func (s *Server) ReturnOrder(ctx context.Context, req *pb.ReturnOrderRequest) (*pb.ReturnOrderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	s.logEvent("ReturnOrder", req)

	cmd := commands.NewReturnOrder(int(req.Order))
	_, err := s.executor.ReturnOrder(cmd)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to execute command: %v", err)
	}

	return &pb.ReturnOrderResponse{}, nil
}
