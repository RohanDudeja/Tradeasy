-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users(
    user_id  varchar(255) not null primary key,
    email_id varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp  not null,
    updated_at timestamp  not null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
