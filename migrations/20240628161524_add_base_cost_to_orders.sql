-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD COLUMN base_cost INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUNN base_cost;
-- +goose StatementEnd
