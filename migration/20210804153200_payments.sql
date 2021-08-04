-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS payments
(
    ID INT AUTO_INCREMENT PRIMARY KEY NOT NULL,
    UserId VARCHAR(255) NOT NULL,
    RazorpayLinkId VARCHAR(255),
    RazorpayLink  VARCHAR(255),
    Amount INT NOT NULL,
    PaymentsType  VARCHAR(10) NOT NULL,
    CurrentBalance INT NOT NULL,
    CreatedAt TIMESTAMP NOT NULL,
    UpdatedAt TIMESTAMP NOT NULL,
    DeletedAt TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE payments;
-- +goose StatementEnd
