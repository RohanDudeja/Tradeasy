-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks(
    id int not null AUTO_INCREMENT primary key,
    stock_name varchar(255) not null,
    open_price int not null,
    stock_ticker_symbol varchar(255) not null,
    ltp int not null,
    high_price int not null,
    low_price int not null,
    previous_day_close int not null,
    percentage_change int not null,
    created_at timestamp  not null,
    updated_at timestamp  not null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
-- +goose StatementEnd