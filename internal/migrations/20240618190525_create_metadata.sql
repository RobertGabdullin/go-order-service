-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS metadata (
    hash TEXT PRIMARY KEY
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS metadata;
-- +goose StatementEnd
