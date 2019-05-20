CREATE TABLE IF NOT EXISTS price_token (
    id SERIAL PRIMARY KEY,
    time_at TIMESTAMP NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    price float(18) NOT NULL
);