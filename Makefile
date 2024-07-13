MIGRATION_DIR=./migrations
DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/orders?sslmode=disable
TEST_DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5433/orders_test?sslmode=disable

PROTO_SRC_DIR := ./proto
PROTO_OUT_DIR := ./pb

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
	go test -tags=unit ./...

test-integration:
	go test -tags=integration ./tests/...

test-db:
	docker-compose -f docker-compose.test.yml up -d

test-migrate-up:
	goose -dir $(MIGRATION_DIR) postgres $(TEST_DB_CONNECTION_STRING) up

clean-test-db:
	docker-compose exec -T postgres psql -U postgres -c "DROP DATABASE IF EXISTS orders_test;"
	docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE orders_test;"

test: test-db clean-test-db test-migrate-up test-unit test-integration

proto:
	protoc --proto_path=$(PROTO_SRC_DIR) --go_out=$(PROTO_OUT_DIR) --go-grpc_out=$(PROTO_OUT_DIR) $(PROTO_SRC_DIR)/*.proto

