package storage

import (
	"database/sql"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

const (
	queryInsertOrder             = `INSERT INTO orders (id, recipient, status, time_limit, delivered_at, returned_at) VALUES ($1, $2, $3, $4, $5, $6)`
	queryUpdateOrder             = `UPDATE orders SET recipient = $1, status = $2, time_limit = $3, delivered_at = $4, returned_at = $5 WHERE id = $6`
	queryDeleteOrder             = `DELETE FROM orders WHERE id = $1`
	querySelectOrderById         = `SELECT id, recipient, status, time_limit, delivered_at, returned_at FROM orders WHERE id = $1`
	querySelectOrdersByRecipient = `SELECT id, recipient, status, time_limit, delivered_at, returned_at FROM orders WHERE recipient = $1`
	querySelectOrdersByStatus    = `SELECT id, recipient, status, time_limit, delivered_at, returned_at FROM orders WHERE status = $1`
	queryUpdateHash              = `UPDATE metadata SET hash = $1`
	queryInsertHash              = `INSERT INTO metadata (hash) VALUES ($1)`
	queryHashExists              = `SELECT COUNT(*) FROM metadata`
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) AddOrder(ord models.Order) error {
	_, err := s.db.Exec(queryInsertOrder, ord.Id, ord.Recipient, ord.Status, ord.Limit, ord.DeliveredAt, ord.ReturnedAt)
	return err
}

func (s *PostgresStorage) UpdateOrder(ord models.Order) error {
	_, err := s.db.Exec(queryUpdateOrder, ord.Recipient, ord.Status, ord.Limit, ord.DeliveredAt, ord.ReturnedAt, ord.Id)
	return err
}

func (s *PostgresStorage) DeleteOrder(id int) error {
	_, err := s.db.Exec(queryDeleteOrder, id)
	return err
}

func (s *PostgresStorage) GetOrderById(id int) (models.Order, error) {
	row := s.db.QueryRow(querySelectOrderById, id)

	var ord models.Order
	err := row.Scan(&ord.Id, &ord.Recipient, &ord.Status, &ord.Limit, &ord.DeliveredAt, &ord.ReturnedAt)
	if err != nil {
		return models.Order{}, err
	}
	return ord, nil
}

func (s *PostgresStorage) GetOrdersByRecipient(recipient int) ([]models.Order, error) {
	rows, err := s.db.Query(querySelectOrdersByRecipient, recipient)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var ord models.Order
		if err := rows.Scan(&ord.Id, &ord.Recipient, &ord.Status, &ord.Limit, &ord.DeliveredAt, &ord.ReturnedAt); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (s *PostgresStorage) GetOrdersByStatus(status string) ([]models.Order, error) {
	rows, err := s.db.Query(querySelectOrdersByStatus, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var ord models.Order
		if err := rows.Scan(&ord.Id, &ord.Recipient, &ord.Status, &ord.Limit, &ord.DeliveredAt, &ord.ReturnedAt); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (s *PostgresStorage) UpdateHash(hash string) error {
	_, err := s.db.Exec(queryUpdateHash, hash)
	return err
}

func (s *PostgresStorage) InsertHash(hash string) error {
	_, err := s.db.Exec(queryInsertHash, hash)
	return err
}

func (s *PostgresStorage) HashExists() (bool, error) {
	var count int
	err := s.db.QueryRow(queryHashExists).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *PostgresStorage) BeginTransaction() (*sql.Tx, error) {
	return s.db.Begin()
}

func (s *PostgresStorage) CommitTransaction(tx *sql.Tx) error {
	return tx.Commit()
}

func (s *PostgresStorage) RollbackTransaction(tx *sql.Tx) error {
	return tx.Rollback()
}
