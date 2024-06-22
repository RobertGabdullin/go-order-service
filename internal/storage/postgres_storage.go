package storage

import (
	"database/sql"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

const (
	queryInsertOrder             = `INSERT INTO orders (id, recipient, status, time_limit, delivered_at, returned_at, hash, weight, total_cost, wrapper) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	queryUpdateOrder             = `UPDATE orders SET recipient = $1, status = $2, time_limit = $3, delivered_at = $4, returned_at = $5, hash = $6 WHERE id = $7`
	queryDeleteOrder             = `DELETE FROM orders WHERE id = $1`
	querySelectOrderById         = `SELECT id, recipient, status, time_limit, delivered_at, returned_at, hash FROM orders WHERE id = $1`
	querySelectOrdersByRecipient = `SELECT id, recipient, status, time_limit, delivered_at, returned_at, hash FROM orders WHERE recipient = $1`
	querySelectOrdersByStatus    = `SELECT id, recipient, status, time_limit, delivered_at, returned_at, hash FROM orders WHERE status = $1 ORDER BY returned_at`
	queryUpdateHash              = `UPDATE orders SET hash = $1 WHERE id = $2`
	querySelectWrapperByType     = `SELECT type, max_weight, markup FROM wrappers WHERE type = $1`
)

type PostgresStorage struct {
	db *sql.DB
}

func New(connStr string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) AddOrder(ord models.Order) error {
	_, err := s.db.Exec(queryInsertOrder, ord.Id, ord.Recipient, ord.Status, ord.Limit, ord.DeliveredAt, ord.ReturnedAt, ord.Hash, ord.Weight, ord.TotalPrice, ord.Wrapper)
	return err
}

func (s *PostgresStorage) UpdateOrder(ord models.Order) error {
	_, err := s.db.Exec(queryUpdateOrder, ord.Recipient, ord.Status, ord.Limit, ord.DeliveredAt, ord.ReturnedAt, ord.Hash, ord.Id)
	return err
}

func (s *PostgresStorage) DeleteOrder(id int) error {
	_, err := s.db.Exec(queryDeleteOrder, id)
	return err
}

func (s *PostgresStorage) GetOrderById(id int) (models.Order, error) {
	row := s.db.QueryRow(querySelectOrderById, id)

	var ord models.Order
	err := row.Scan(&ord.Id, &ord.Recipient, &ord.Status, &ord.Limit, &ord.DeliveredAt, &ord.ReturnedAt, &ord.Hash)
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
		if err := rows.Scan(&ord.Id, &ord.Recipient, &ord.Status, &ord.Limit, &ord.DeliveredAt, &ord.ReturnedAt, &ord.Hash); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (s *PostgresStorage) GetPaginatedOrdersByStatus(status string, offset, limit int) ([]models.Order, error) {

	var rows *sql.Rows
	var err error
	query := querySelectOrdersByStatus
	if limit > 0 {
		query = query + " LIMIT $2 OFFSET $3"
		rows, err = s.db.Query(query, status, limit, offset)
	} else {
		query = query + " OFFSET $2"
		rows, err = s.db.Query(query, status, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var ord models.Order
		if err := rows.Scan(&ord.Id, &ord.Recipient, &ord.Status, &ord.Limit, &ord.DeliveredAt, &ord.ReturnedAt, &ord.Hash); err != nil {
			return nil, err
		}
		orders = append(orders, ord)
	}
	return orders, nil
}

func (s *PostgresStorage) UpdateHash(id int, hash string) error {
	_, err := s.db.Exec(queryUpdateHash, hash, id)
	return err
}

func (s *PostgresStorage) GetWrapperByType(givenType string) ([]models.Wrapper, error) {
	rows, err := s.db.Query(querySelectWrapperByType, givenType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wrappers []models.Wrapper
	for rows.Next() {
		var wrapper models.Wrapper
		if err := rows.Scan(&wrapper.Type, &wrapper.MaxWeight, &wrapper.Markup); err != nil {
			return nil, err
		}
		wrappers = append(wrappers, wrapper)
	}
	return wrappers, nil
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
