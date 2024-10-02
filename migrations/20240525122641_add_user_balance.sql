-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD balance DOUBLE PRECISION NOT NULL DEFAULT 0,
    ADD withdrawn DOUBLE PRECISION NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP balance,
    DROP withdrawn;
-- +goose StatementEnd
