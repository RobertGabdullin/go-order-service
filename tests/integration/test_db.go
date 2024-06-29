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

	schema := `
	CREATE TABLE IF NOT EXISTS orders (
		id PRIMARY KEY,
		recipient INT,
		status TEXT,
		time_limit TIMESTAMP,
		delivered_at TIMESTAMP,
		returned_at TIMESTAMP,
		hash TEXT,
		weight INT,
		base_price INT,
		wrapper TEXT
	);

	CREATE TABLE IF NOT EXISTS wrappers (
		id SERIAL PRIMARY KEY,
		type TEXT,
		max_weight INT,
		markup INT
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func teardownDB(db *sql.DB) {
	db.Exec("DROP TABLE IF EXISTS orders;")
	db.Exec("DROP TABLE IF EXISTS wrappers;")
	db.Close()
}
