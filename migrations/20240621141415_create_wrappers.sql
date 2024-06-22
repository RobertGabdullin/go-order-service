-- +goose Up
-- +goose StatementBegin
CREATE TABLE wrappers (
    id INT PRIMARY KEY,
    type TEXT NOT NULL,
    max_weight INT,
    markup INT NOT NULL
);

INSERT INTO wrappers (id, type, max_weight, markup) VALUES
(1, 'pack', 10, 5),
(2, 'box', 30, 20),
(3, 'film', NULL, 1), 
(4, 'none', NULL, 0);  

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wrappers;
-- +goose StatementEnd
