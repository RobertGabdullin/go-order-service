//go:build unit

package api

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	pb "gitlab.ozon.dev/r_gabdullin/homework-1/pb"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestAcceptOrder(t *testing.T) {
	t.Parallel()

	mockStorageService := new(tests.MockStorageService)
	s := NewServer(mockStorageService, nil)

	mockStorageService.On("AddOrder", mock.AnythingOfType("models.Order")).Return(nil)

	mockStorageService.On("GetWrapper", "wrapper").Return(models.Wrapper{}, nil)

	req := &pb.AcceptOrderRequest{
		User:      1,
		Order:     1,
		Weight:    1,
		BasePrice: 1,
		Expire:    "2024-12-31T23",
		Wrapper:   "wrapper",
	}

	resp, err := s.AcceptOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, &emptypb.Empty{}, resp)
	mockStorageService.AssertExpectations(t)
}

func TestAcceptReturn(t *testing.T) {
	t.Parallel()
	mockStorageService := new(tests.MockStorageService)

	s := NewServer(mockStorageService, nil)

	mockStorageService.On("ChangeStatus", 1, "returned", mock.AnythingOfType("string")).Return(nil)
	mockStorageService.On("FindOrders", []int{1}).Return([]models.Order{
		{
			Id:          1,
			Recipient:   1,
			Status:      "delivered",
			DeliveredAt: time.Now().Add(-24 * time.Hour),
		},
	}, nil)

	req := &pb.AcceptReturnRequest{
		User:  1,
		Order: 1,
	}

	resp, err := s.AcceptReturn(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, &emptypb.Empty{}, resp)
	mockStorageService.AssertExpectations(t)
}

func TestDeliverOrder(t *testing.T) {
	t.Parallel()

	mockStorageService := new(tests.MockStorageService)
	s := NewServer(mockStorageService, nil)

	mockStorageService.On("FindOrders", []int{1, 2}).Return([]models.Order{
		{
			Id:        1,
			Recipient: 1,
			Status:    "alive",
			Expire:    time.Now().Add(24 * time.Hour),
		},
		{
			Id:        2,
			Recipient: 1,
			Status:    "alive",
			Expire:    time.Now().Add(24 * time.Hour),
		},
	}, nil)
	mockStorageService.On("ChangeStatus", 1, "delivered", mock.AnythingOfType("string")).Return(nil)
	mockStorageService.On("ChangeStatus", 2, "delivered", mock.AnythingOfType("string")).Return(nil)

	req := &pb.DeliverOrderRequest{
		Orders: []int32{1, 2},
	}

	resp, err := s.DeliverOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, &emptypb.Empty{}, resp)
	mockStorageService.AssertExpectations(t)
}

func TestGetOrders(t *testing.T) {
	t.Parallel()

	mockStorageService := new(tests.MockStorageService)
	s := NewServer(mockStorageService, nil)

	orders := []models.Order{
		{Id: 1, Recipient: 2, Expire: time.Now(), Status: "alive"},
	}

	mockStorageService.On("ListOrders", 2).Return(orders, nil)

	req := &pb.GetOrdersRequest{
		User:  2,
		Count: 1,
	}

	resp, err := s.GetOrders(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, resp.Orders, 1)
	assert.Equal(t, int32(1), resp.Orders[0].Id)
	assert.Equal(t, int32(2), resp.Orders[0].Recipient)
	assert.Equal(t, "alive", resp.Orders[0].Status)
	mockStorageService.AssertExpectations(t)
}

func TestGetReturns(t *testing.T) {
	t.Parallel()

	mockStorageService := new(tests.MockStorageService)
	s := NewServer(mockStorageService, nil)

	returns := []models.Order{
		{Id: 1, Recipient: 2, Expire: time.Now(), ReturnedAt: time.Now()},
	}

	mockStorageService.On("GetReturns", 1, 1).Return(returns, nil)

	req := &pb.GetReturnsRequest{
		Offset: 1,
		Limit:  1,
	}

	resp, err := s.GetReturns(context.Background(), req)

	assert.NoError(t, err)
	assert.Len(t, resp.Returns, 1)
	assert.Equal(t, int32(1), resp.Returns[0].Id)
	assert.Equal(t, int32(2), resp.Returns[0].Recipient)
	mockStorageService.AssertExpectations(t)
}

func TestReturnOrder(t *testing.T) {
	t.Parallel()

	mockStorageService := new(tests.MockStorageService)
	s := NewServer(mockStorageService, nil)

	mockStorageService.On("FindOrders", []int{1}).Return([]models.Order{
		{
			Id:     1,
			Status: "alive",
			Expire: time.Now().Add(-24 * time.Hour),
		},
	}, nil)
	mockStorageService.On("DeleteOrder", 1).Return(nil)

	req := &pb.ReturnOrderRequest{
		Order: 1,
	}

	resp, err := s.ReturnOrder(context.Background(), req)

	assert.NoError(t, err)
	assert.Equal(t, &emptypb.Empty{}, resp)
	mockStorageService.AssertExpectations(t)
}
