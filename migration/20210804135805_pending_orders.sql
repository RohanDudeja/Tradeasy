-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pending_orders(
       userid varchar(255) not null ,
       order_id varchar(255) not null primary key ,
       stock_name varchar(255) not null ,
       order_type varchar(255) not null ,
       book_type varchar(255) not null ,
       limit_price int not null ,
       quantity int not null ,
       order_price int not null ,
       status varchar(255) not null ,
       created_at timestamp default current_timestamp,
       updated_at timestamp default null,
       deleted_at timestamp default null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pending_orders;
-- +goose StatementEnd
