MIGRATION_DIR=./migrations
DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5432/orders?sslmode=disable
TEST_DB_CONNECTION_STRING=postgres://postgres:postgres@localhost:5433/orders_test?sslmode=disable

PROTO_DIR=./proto
PROTO_OUT_DIR=./pb
DEP_DIR=./dep
PROTO_FILES=./proto/*.proto

GOOGLEAPIS_REPO=https://github.com/googleapis/googleapis.git
GRPC_GATEWAY_REPO=https://github.com/grpc-ecosystem/grpc-gateway.git
PROTOC_GEN_VALIDATE_REPO=https://github.com/envoyproxy/protoc-gen-validate.git

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

test-environment:
	docker-compose -f docker-compose.test.yml up -d

test-migrate-up:
	goose -dir $(MIGRATION_DIR) postgres $(TEST_DB_CONNECTION_STRING) up

clean-test-db:
	docker-compose exec -T postgres psql -U postgres -c "DROP DATABASE IF EXISTS orders_test;"
	docker-compose exec -T postgres psql -U postgres -c "CREATE DATABASE orders_test;"

test: test-db clean-test-db test-migrate-up test-unit test-integration

clone-deps:
	mkdir -p $(DEP_DIR)
	git clone $(GOOGLEAPIS_REPO) $(DEP_DIR)/googleapis
	git clone $(GRPC_GATEWAY_REPO) $(DEP_DIR)/grpc-gateway
	git clone $(PROTOC_GEN_VALIDATE_REPO) $(DEP_DIR)/protoc-gen-validate

install-deps: clone-deps
	go mod tidy
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	go install github.com/envoyproxy/protoc-gen-validate@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

generate:
	protoc -I $(DEP_DIR)/googleapis \
	       -I $(DEP_DIR)/grpc-gateway \
	       -I $(DEP_DIR)/protoc-gen-validate \
	       --proto_path=$(PROTO_DIR) \
	       --go_out=$(PROTO_OUT_DIR) \
	       --go-grpc_out=$(PROTO_OUT_DIR) \
	       --grpc-gateway_out=logtostderr=true:$(PROTO_OUT_DIR) \
	       --validate_out="lang=go:$(PROTO_OUT_DIR)" \
	       $(PROTO_FILES)

