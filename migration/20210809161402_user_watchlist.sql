-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_watchlist(
    id int not null AUTO_INCREMENT primary key,
    user_id varchar(255) not null,
    watchlist_id varchar(255) not null,
    stock_name varchar(255) not null,
    created_at timestamp  not null,
    updated_at timestamp  not null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_watchlist;
-- +goose StatementEnd
