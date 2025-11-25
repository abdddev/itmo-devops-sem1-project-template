-- +goose Up
CREATE TABLE IF NOT EXISTS prices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    create_date TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_prices_unique
    ON prices (name, category, price, create_date);