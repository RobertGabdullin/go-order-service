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
		id INT PRIMARY KEY,
		recipient INT,
		status TEXT,
		time_limit TIMESTAMPTZ,
		delivered_at TIMESTAMPTZ,
		returned_at TIMESTAMPTZ,
		hash TEXT,
		weight INT,
		base_cost INT,
		wrapper TEXT
	);

	CREATE TABLE IF NOT EXISTS wrappers (
		id INT PRIMARY KEY,
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
