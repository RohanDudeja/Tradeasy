-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trading_account(
    user_id  varchar(255) not null,
    id int not null AUTO_INCREMENT primary key,
    pan_card_no varchar(255) not null,
    bank_acc_no varchar(255) not null,
    trading_acc_no varchar(255) not null,
    balance bigint not null,
    created_at timestamp  not null,
    updated_at timestamp  not null,
    deleted_at timestamp  default null
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE trading_account;
-- +goose StatementEnd
