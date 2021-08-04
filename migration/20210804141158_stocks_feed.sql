-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks_feed(
        id int not null AUTO_INCREMENT primary key ,
        stock_name varchar(255) not null ,
        open int not null ,
        ltp int not null ,
        high int not null ,
        low int not null ,
        traded_at timestamp default null,
        created_at timestamp default current_timestamp,
        deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks_feed;
-- +goose StatementEnd
