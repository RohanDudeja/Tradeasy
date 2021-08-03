-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS BUY_ORDER_BOOK(
    order_id varchar not null,
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
DROP TABLE BUY_ORDER_BOOK;
-- +goose StatementEnd
