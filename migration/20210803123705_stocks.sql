-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS STOCKS(
    id int not null unique,
    stock_name varchar,
    open_price int,
    stock_ticker_symbol varchar primary key,
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
DROP TABLE STOCKS;
-- +goose StatementEnd