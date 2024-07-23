package executor_test

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/executor"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
	"gitlab.ozon.dev/r_gabdullin/homework-1/tests"
)

func TestAcceptOrder_Success(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.AcceptOrder{
		Order:     1,
		Recipient: 1,
		Expire:    time.Now().Add(24 * time.Hour),
		Weight:    10,
		BasePrice: 100,
		Wrapper:   "standard",
	}

	mockService.On("GetWrapper", "standard").Return(models.Wrapper{
		Type:      "standard",
		MaxWeight: sql.NullInt64{Int64: 50, Valid: true},
	}, nil)
	mockService.On("AddOrder", mock.AnythingOfType("models.Order")).Return(nil)

	orders, err := exec.AcceptOrder(args)

	assert.NoError(t, err)
	assert.Len(t, orders, 0)
	mockService.AssertExpectations(t)
}

func TestAcceptOrder_InvalidWeight(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.AcceptOrder{
		Order:     1,
		Recipient: 1,
		Expire:    time.Now().Add(24 * time.Hour),
		Weight:    100,
		BasePrice: 100,
		Wrapper:   "standard",
	}

	mockService.On("GetWrapper", "standard").Return(models.Wrapper{
		Type:      "standard",
		MaxWeight: sql.NullInt64{Int64: 50, Valid: true},
	}, nil)

	orders, err := exec.AcceptOrder(args)

	assert.Error(t, err)
	assert.Equal(t, "order weight exceeds the maximum limit for the chosen wrapper", err.Error())
	assert.Nil(t, orders)
	mockService.AssertExpectations(t)
}

func TestReturnOrder_Success(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.ReturnOrder{
		Order: 1,
	}

	mockService.On("FindOrders", []int{1}).Return([]models.Order{
		{Id: 1, Status: "alive", Expire: time.Now().Add(-24 * time.Hour)},
	}, nil)
	mockService.On("DeleteOrder", 1).Return(nil)

	orders, err := exec.ReturnOrder(args)

	assert.NoError(t, err)
	assert.Len(t, orders, 0)
	mockService.AssertExpectations(t)
}

func TestReturnOrder_NotExist(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.ReturnOrder{
		Order: 1,
	}

	mockService.On("FindOrders", []int{1}).Return(nil, errors.New("such order does not exist"))

	orders, err := exec.ReturnOrder(args)

	assert.Error(t, err)
	assert.Equal(t, "such order does not exist", err.Error())
	assert.Nil(t, orders)
	mockService.AssertExpectations(t)
}

func TestAcceptReturn_Success(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.AcceptReturn{
		Order: 1,
	}

	mockService.On("FindOrders", []int{1}).Return([]models.Order{
		{Id: 1, Status: "delivered", DeliveredAt: time.Now().Add(-1 * time.Hour)},
	}, nil)
	mockService.On("ChangeStatus", 1, "returned", mock.AnythingOfType("string")).Return(nil)

	orders, err := exec.AcceptReturn(args)

	assert.NoError(t, err)
	assert.Len(t, orders, 0)
	mockService.AssertExpectations(t)
}

func TestDeliverOrder_Success(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.DeliverOrder{
		Ords: []int{1, 2, 3},
	}

	mockService.On("FindOrders", []int{1, 2, 3}).Return([]models.Order{
		{Id: 1, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
		{Id: 2, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
		{Id: 3, Recipient: 1, Status: "alive", Expire: time.Now().Add(24 * time.Hour)},
	}, nil)
	mockService.On("ChangeStatus", 1, "delivered", mock.AnythingOfType("string")).Return(nil)
	mockService.On("ChangeStatus", 2, "delivered", mock.AnythingOfType("string")).Return(nil)
	mockService.On("ChangeStatus", 3, "delivered", mock.AnythingOfType("string")).Return(nil)

	orders, err := exec.DeliverOrder(args)

	assert.NoError(t, err)
	assert.Len(t, orders, 0)
	mockService.AssertExpectations(t)
}

func TestGetOrders_Success(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.GetOrders{
		User:  1,
		Count: 5,
	}

	expectedOrders := []models.Order{
		{Id: 1, Recipient: 1, Expire: time.Now(), Status: "alive"},
		{Id: 2, Recipient: 1, Expire: time.Now(), Status: "alive"},
	}

	mockService.On("ListOrders", 1).Return(expectedOrders, nil)

	orders, err := exec.GetOrders(args)

	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	mockService.AssertExpectations(t)
}

func TestGetReturns_Success(t *testing.T) {
	mockService := new(tests.MockStorageService)
	exec := executor.NewOrderCommandExecutor(mockService)

	args := commands.GetReturns{
		Offset: 0,
		Limit:  5,
	}

	expectedReturns := []models.Order{
		{Id: 1, Recipient: 1, Expire: time.Now(), ReturnedAt: time.Now()},
		{Id: 2, Recipient: 2, Expire: time.Now(), ReturnedAt: time.Now()},
	}

	mockService.On("GetReturns", 0, 5).Return(expectedReturns, nil)

	orders, err := exec.GetReturns(args)

	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	mockService.AssertExpectations(t)
}
