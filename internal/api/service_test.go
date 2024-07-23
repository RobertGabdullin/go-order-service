package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAcceptOrder_Success(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.AcceptOrderRequest{
		User:      1,
		Order:     1,
		Weight:    10,
		BasePrice: 100,
		Expire:    "2025-12-31T12",
		Wrapper:   "pack",
	}

	mockExecutor.On("AcceptOrder", mock.Anything).Return(nil, nil)

	resp, err := server.AcceptOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockExecutor.AssertExpectations(t)
}

func TestAcceptOrder_InvalidRequest(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.AcceptOrderRequest{
		User:      1,
		Order:     1,
		Weight:    10,
		BasePrice: 100,
		Expire:    "invalid-date",
		Wrapper:   "pack",
	}

	resp, err := server.AcceptOrder(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestAcceptReturn_Success(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.AcceptReturnRequest{
		User:  1,
		Order: 1,
	}

	cmd := commands.NewAcceptReturn(int(req.User), int(req.Order))

	mockExecutor.On("AcceptReturn", cmd).Return(nil, nil)

	resp, err := server.AcceptReturn(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockExecutor.AssertExpectations(t)
}

func TestAcceptReturn_InvalidRequest(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.AcceptReturnRequest{
		User:  1,
		Order: -1,
	}

	resp, err := server.AcceptReturn(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestDeliverOrder_Success(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.DeliverOrderRequest{
		Orders: []int32{1, 2, 3},
	}

	orders := []int{1, 2, 3}
	cmd := commands.NewDeliverOrder(orders)

	mockExecutor.On("DeliverOrder", cmd).Return(nil, nil)

	resp, err := server.DeliverOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockExecutor.AssertExpectations(t)
}

func TestDeliverOrder_InvalidRequest(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.DeliverOrderRequest{
		Orders: []int32{-1},
	}

	resp, err := server.DeliverOrder(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetOrders_Success(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.GetOrdersRequest{
		User:  1,
		Count: 5,
	}

	cmd := commands.NewGetOrders(int(req.User), int(req.Count))
	expectedOrders := []models.Order{
		{Id: 1, Recipient: 1, Expire: time.Now(), Status: "alive"},
		{Id: 2, Recipient: 1, Expire: time.Now(), Status: "alive"},
	}

	mockExecutor.On("GetOrders", cmd).Return(expectedOrders, nil)

	resp, err := server.GetOrders(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Orders, 2)
	mockExecutor.AssertExpectations(t)
}

func TestGetOrders_InvalidRequest(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.GetOrdersRequest{
		User:  -1,
		Count: 5,
	}

	resp, err := server.GetOrders(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestGetReturns_Success(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.GetReturnsRequest{
		Offset: 0,
		Limit:  5,
	}

	cmd := commands.NewGetReturns(int(req.Offset), int(req.Limit))
	expectedReturns := []models.Order{
		{Id: 1, Recipient: 1, Expire: time.Now(), ReturnedAt: time.Now(), Status: "returned"},
		{Id: 2, Recipient: 2, Expire: time.Now(), ReturnedAt: time.Now(), Status: "returned"},
	}

	mockExecutor.On("GetReturns", cmd).Return(expectedReturns, nil)

	resp, err := server.GetReturns(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Len(t, resp.Returns, 2)
	mockExecutor.AssertExpectations(t)
}

func TestGetReturns_InvalidRequest(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.GetReturnsRequest{
		Offset: -1,
		Limit:  5,
	}

	resp, err := server.GetReturns(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestReturnOrder_Success(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.ReturnOrderRequest{
		Order: 1,
	}

	cmd := commands.NewReturnOrder(int(req.Order))

	mockExecutor.On("ReturnOrder", cmd).Return(nil, nil)

	resp, err := server.ReturnOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	mockExecutor.AssertExpectations(t)
}

func TestReturnOrder_InvalidRequest(t *testing.T) {
	mockExecutor := new(tests.CommandExecutorMock)
	server := NewServer(mockExecutor, nil)

	req := &pb.ReturnOrderRequest{
		Order: -1,
	}

	resp, err := server.ReturnOrder(context.Background(), req)

	assert.Nil(t, resp)
	assert.Error(t, err)
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}
