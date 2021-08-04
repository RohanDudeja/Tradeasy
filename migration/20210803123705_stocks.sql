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
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
-- +goose StatementEnd