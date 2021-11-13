-- +goose Up
-- +goose StatementBegin
CREATE TABLE jobstore (
     id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
     type VARCHAR(50),
     next_time_execute TIMESTAMP,
     parameters json,
     count INT,
     executed BOOL DEFAULT FALSE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE jobstore;
-- +goose StatementEnd
