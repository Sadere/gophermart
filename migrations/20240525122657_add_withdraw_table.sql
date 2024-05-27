-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS withdrawals (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    number varchar(255) NOT NULL,
    created_at timestamp NOT NULL,
    amount DOUBLE PRECISION NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE withdrawals; 
-- +goose StatementEnd
