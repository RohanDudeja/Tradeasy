-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS stocks(
    id int not null unique AUTO_INCREMENT primary key,
    stock_name varchar,
    open_price int,
    stock_ticker_symbol varchar ,
    ltp int,
    high_price int,
    low_price int,
    previous_day_close int,
    percentage_change int,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stocks;
-- +goose StatementEnd