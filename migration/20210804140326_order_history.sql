-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_history(
        userid varchar(255) not null ,
        order_id varchar(255) not null ,
        id int not null AUTO_INCREMENT primary key ,
        stock_name varchar(255) not null ,
        quantity int not null ,
        buy_price int not null ,
        sell_price int not null ,
        commission_fee int not null ,
        bought_at timestamp not null ,
        sold_at timestamp not null ,
        created_at timestamp not null ,
        updated_at timestamp not null ,
        deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE order_history;
-- +goose StatementEnd
