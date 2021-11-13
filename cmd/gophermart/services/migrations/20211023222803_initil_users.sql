-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    login VARCHAR(50) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    balance FLOAT DEFAULT 0,
    spend FLOAT DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
