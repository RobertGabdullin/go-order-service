package tests

import (
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type MockStorageService struct {
	mock.Mock
}

func (m *MockStorageService) AddOrder(order models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockStorageService) GetWrapper(name string) (models.Wrapper, error) {
	args := m.Called(name)
	return args.Get(0).(models.Wrapper), args.Error(1)
}

func (m *MockStorageService) FindOrders(ids []int) ([]models.Order, error) {
	args := m.Called(ids)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *MockStorageService) ChangeStatus(id int, status, hash string) error {
	args := m.Called(id, status, hash)
	return args.Error(0)
}

func (m *MockStorageService) ListOrders(userID int) ([]models.Order, error) {
	args := m.Called(userID)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *MockStorageService) GetReturns(offset, limit int) ([]models.Order, error) {
	args := m.Called(offset, limit)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *MockStorageService) DeleteOrder(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
