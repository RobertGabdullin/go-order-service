package integration

import (
	"database/sql"
)

var (
	dbConnStr = "postgres://postgres:postgres@localhost:5433/orders_test?sslmode=disable"
)

func setupDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func teardownDB(db *sql.DB) {
	db.Close()
}
