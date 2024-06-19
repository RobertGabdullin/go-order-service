package service

import (
	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type StorageService interface {
	AddOrder(models.Order) error
	ChangeStatus(id int, status string) error
	FindOrders(ids []int) ([]models.Order, error)
	ListOrders(recipient int) ([]models.Order, error)
	GetReturns() ([]models.Order, error)
	UpdateHash(hash string) error
	DeleteOrder(id int) error
}
