-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN wrapper TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUNN wrapper;
-- +goose StatementEnd
