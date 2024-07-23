package tests

import (
	"github.com/stretchr/testify/mock"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/commands"
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type CommandExecutorMock struct {
	mock.Mock
}

func (m *CommandExecutorMock) AcceptOrder(cmd commands.AcceptOrder) ([]models.Order, error) {
	args := m.Called(cmd)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *CommandExecutorMock) AcceptReturn(cmd commands.AcceptReturn) ([]models.Order, error) {
	args := m.Called(cmd)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *CommandExecutorMock) DeliverOrder(cmd commands.DeliverOrder) ([]models.Order, error) {
	args := m.Called(cmd)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *CommandExecutorMock) GetOrders(cmd commands.GetOrders) ([]models.Order, error) {
	args := m.Called(cmd)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *CommandExecutorMock) GetReturns(cmd commands.GetReturns) ([]models.Order, error) {
	args := m.Called(cmd)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}

func (m *CommandExecutorMock) ReturnOrder(cmd commands.ReturnOrder) ([]models.Order, error) {
	args := m.Called(cmd)
	var orders []models.Order
	if args.Get(0) != nil {
		orders = args.Get(0).([]models.Order)
	}
	return orders, args.Error(1)
}
