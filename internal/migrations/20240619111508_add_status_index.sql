-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_status ON orders (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_status;
-- +goose StatementEnd
