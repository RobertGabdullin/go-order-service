package storage

import (
	"database/sql"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

type Storage interface {
	AddOrder(models.Order) error
	UpdateOrder(models.Order) error
	DeleteOrder(id int) error
	GetOrderById(id int) (models.Order, error)
	GetOrdersByRecipient(recipient int) ([]models.Order, error)
	GetOrdersByStatus(status string) ([]models.Order, error)
	UpdateHash(hash string) error
	InsertHash(hash string) error
	HashExists() (bool, error)
}

type TransactionalStorage interface {
	Storage
	BeginTransaction() (*sql.Tx, error)
	CommitTransaction(tx *sql.Tx) error
	RollbackTransaction(tx *sql.Tx) error
}
