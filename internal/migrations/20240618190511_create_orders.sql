-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    id INT PRIMARY KEY,
    recipient INT NOT NULL,
    status TEXT NOT NULL,
    time_limit TIMESTAMP,
    delivered_at TIMESTAMP,
    returned_at TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
