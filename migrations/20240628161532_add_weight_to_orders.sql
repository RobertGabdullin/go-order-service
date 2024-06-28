-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN weight INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUNN weight;
-- +goose StatementEnd
