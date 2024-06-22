MIGRATION_DIR=./migrations

DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/orders?sslmode=disable

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
