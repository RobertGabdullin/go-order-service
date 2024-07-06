-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    id INT PRIMARY KEY,
    recipient INT NOT NULL,
    status TEXT NOT NULL,
    time_limit TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    returned_at TIMESTAMPTZ,
    hash TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
