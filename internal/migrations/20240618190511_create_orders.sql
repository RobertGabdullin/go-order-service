-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id INT PRIMARY KEY,
    recipient INT NOT NULL,
    status TEXT NOT NULL,
    time_limit TIMESTAMP,
    delivered_at TIMESTAMP,
    returned_at TIMESTAMP,
    hash TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
