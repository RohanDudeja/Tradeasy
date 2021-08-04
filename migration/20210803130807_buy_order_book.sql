-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS buy_order_book(
    id int not null AUTO_INCREMENT primary key,
    order_id  varchar(255) not null,
    foreign key (order_id) references pending_orders(order_id),
    stock_ticker_symbol varchar(255),
    order_quantity int,
    order_status varchar(255),
    order_price int,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE buy_order_book;
-- +goose StatementEnd
