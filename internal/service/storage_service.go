package service

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type StorageService interface {
	AddOrder(models.Order) error
	ChangeStatus(id int, status, hash string) error
	FindOrders(ids []int) ([]models.Order, error)
	ListOrders(recipient int) ([]models.Order, error)
	GetReturns(offset, limit int) ([]models.Order, error)
	DeleteOrder(id int) error
	GetWrapper(givenType string) (models.Wrapper, error)
}
