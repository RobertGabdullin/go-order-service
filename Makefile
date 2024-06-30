MIGRATION_DIR=./migrations

DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/orders?sslmode=disable

TEST_DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5433/orders_test?sslmode=disable

install-goose:
	go get -u github.com/pressly/goose/cmd/goose

init-migrations:
	mkdir -p $(MIGRATION_DIR)

new-migration:
	goose -dir $(MIGRATION_DIR) create $(name) sql

migrate-up:
	goose -dir $(MIGRATION_DIR) postgres $(DB_CONNECTION_STRING) up

migrate-down:
	goose -dir $(MIGRATION_DIR) postgres $(DB_CONNECTION_STRING) down

migrate-status:
	goose -dir $(MIGRATION_DIR) postgres $(DB_CONNECTION_STRING) status

test-unit:
	go test -tags=unit .\tests\unit

test-integration:
	go test -tags=integration .\tests\integration

test-db:
	docker-compose -f docker-compose.test.yml up -d

clean-test-db:
	psql -U postgres -c "DROP DATABASE IF EXISTS orders_test; CREATE DATABASE orders_test;"

test: clean-test-db test-unit test-integration




