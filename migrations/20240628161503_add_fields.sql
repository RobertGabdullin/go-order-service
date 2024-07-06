-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders 
    ADD COLUMN wrapper TEXT,
    ADD COLUMN base_cost INT,
    ADD COLUMN weight INT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders 
    DROP COLUNN wrapper,
    DROP COLUMN base_cost,
    DROP COLUMN weight;
-- +goose StatementEnd
