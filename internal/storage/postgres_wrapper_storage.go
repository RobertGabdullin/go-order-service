package storage

import (
	"database/sql"

	"gitlab.ozon.dev/r_gabdullin/homework-1/internal/models"
)

const (
	querySelectWrapperByType = `SELECT id, type, max_weight, markup FROM wrappers WHERE type = $1`
)

type PostgresWrapperStorage struct {
	db *sql.DB
}

func NewWrapperStorage(connStr string) (*PostgresWrapperStorage, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresWrapperStorage{db: db}, nil
}

func (s *PostgresWrapperStorage) GetWrapperByType(givenType string) ([]models.Wrapper, error) {
	rows, err := s.db.Query(querySelectWrapperByType, givenType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wrappers []models.Wrapper
	for rows.Next() {
		var wrapper models.Wrapper
		if err := rows.Scan(&wrapper.Id, &wrapper.Type, &wrapper.MaxWeight, &wrapper.Markup); err != nil {
			return nil, err
		}
		wrappers = append(wrappers, wrapper)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wrappers, nil
}
