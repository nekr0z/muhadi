CREATE TABLE IF NOT EXISTS users(
    username VARCHAR (50) PRIMARY KEY,
    hash VARCHAR (64)
);
CREATE TABLE IF NOT EXISTS orders(
    id BIGINT PRIMARY KEY,
    username VARCHAR (50),
    status SMALLINT,
    accrual DOUBLE PRECISION,
    created_at TIMESTAMPTZ,
    last_updated_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS orders_username_idx ON orders(username);
CREATE INDEX IF NOT EXISTS orders_status_idx ON orders(status);
CREATE TABLE IF NOT EXISTS withdrawals(
    id BIGINT PRIMARY KEY,
    username VARCHAR (50),
    amount DOUBLE PRECISION,
    at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS withdrawals_username_idx ON withdrawals(username);
