-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
   id uuid DEFAULT uuid_generate_v4 (),
   user_id uuid REFERENCES users(id) ON DELETE CASCADE ,
   number VARCHAR(50) NOT NULL UNIQUE,
   status VARCHAR(50),
   uploaded_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   accrual INT DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
