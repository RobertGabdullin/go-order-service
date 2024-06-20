-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_orders_recipient ON orders (recipient);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_orders_recipient;
-- +goose StatementEnd
