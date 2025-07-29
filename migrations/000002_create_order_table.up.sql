CREATE TABLE IF NOT EXISTS orders
(
    id BIGSERIAL PRIMARY KEY,
    uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id BIGSERIAL REFERENCES users(id),
    status INT NULL,
    accrual DECIMAL(10, 2)
);

CREATE INDEX IF NOT EXISTS orders_user_id_ix ON orders(user_id ASC);

CREATE INDEX IF NOT EXISTS orders_status_ix ON orders(status ASC);