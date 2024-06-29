package unit

import (
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type mockStorageService struct {
	mock.Mock
}

func (m *mockStorageService) AddOrder(order models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *mockStorageService) GetWrapper(name string) (models.Wrapper, error) {
	args := m.Called(name)
	return args.Get(0).(models.Wrapper), args.Error(1)
}

func (m *mockStorageService) FindOrders(ids []int) ([]models.Order, error) {
	args := m.Called(ids)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *mockStorageService) ChangeStatus(id int, status, hash string) error {
	args := m.Called(id, status, hash)
	return args.Error(0)
}

func (m *mockStorageService) ListOrders(userID int) ([]models.Order, error) {
	args := m.Called(userID)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *mockStorageService) GetReturns(offset, limit int) ([]models.Order, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.Order), args.Error(1)
}

func (m *mockStorageService) DeleteOrder(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
