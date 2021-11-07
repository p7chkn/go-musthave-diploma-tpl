-- +goose Up
-- +goose StatementBegin
CREATE TABLE withdrawals (
    id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
    user_id uuid REFERENCES users(id) ON DELETE CASCADE ,
    order_number VARCHAR (50) NOT NULL UNIQUE ,
    status VARCHAR(50) DEFAULT 'NEW',
    processed_at TIMESTAMP,
    sum INT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE withdrawals;
-- +goose StatementEnd
