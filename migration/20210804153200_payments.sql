-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments
(
    id INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    razorpay_link_id VARCHAR(255),
    razorpay_link  VARCHAR(255),
    amount BIGINT NOT NULL,
    payments_type  VARCHAR(10) NOT NULL,
    current_balance BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
-- +goose StatementEnd
