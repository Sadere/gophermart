-- +goose Up
-- +goose StatementBegin
CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    created_at timestamp NOT NULL,
    number varchar(255) NOT NULL,
    status order_status DEFAULT 'NEW',
    accrual DOUBLE PRECISION NULL
);
CREATE INDEX number_idx ON orders (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
DROP TYPE IF EXISTS order_status;
-- +goose StatementEnd
