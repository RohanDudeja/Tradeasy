-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sell_order_book(
    order_id varchar not null primary key,
    stock_ticker_symbol varchar,
    order_quantity int,
    order_status varchar,
    order_price varchar,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sell_order_book;
-- +goose StatementEnd
