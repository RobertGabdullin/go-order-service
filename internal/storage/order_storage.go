package storage

import (
	"database/sql"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type OrderStorage interface {
	AddOrder(models.Order) error
	UpdateOrder(models.Order) error
	DeleteOrder(id int) error
	GetOrderById(id int) (models.Order, error)
	GetOrdersByRecipient(recipient int) ([]models.Order, error)
	GetPaginatedOrdersByStatus(status string, offset, limit int) ([]models.Order, error)
	UpdateHash(id int, hash string) error
}

type TransactionalOrderStorage interface {
	OrderStorage
	BeginTransaction() (*sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
	RollbackTransaction(tx *sql.Tx) error
}
