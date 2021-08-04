-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sell_order_book(
    id int not null AUTO_INCREMENT primary key,
    order_id varchar(255) not null,
    stock_ticker_symbol varchar(255) not null,
    order_quantity int not null,
    order_status varchar(255) default 'pending',
    order_price int not null,
    created_at timestamp not null,
    updated_at timestamp not null,
    deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sell_order_book;
-- +goose StatementEnd
