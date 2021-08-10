-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS watchlist(
    id int not null AUTO_INCREMENT primary key,
    name varchar(255) not null,
    created_at timestamp  not null,
    updated_at timestamp  not null,
    deleted_at timestamp default null
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE watchlist;
-- +goose StatementEnd

