CREATE TABLE IF NOT EXISTS transactions
(
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id BIGSERIAL REFERENCES users(id),
    type INT,
    amount DECIMAL(10, 2),
    order_id BIGSERIAL REFERENCES orders(id)
);

CREATE INDEX IF NOT EXISTS transactions_user_id_ix ON transactions(user_id ASC);

CREATE OR REPLACE VIEW bonus_balances AS
SELECT
    user_id,
    SUM(CASE WHEN type = 0 THEN amount ELSE 0 END) AS accrued,
    SUM(CASE WHEN type = 1 THEN amount ELSE 0 END) AS withdrawn,
    SUM(CASE WHEN type = 0 THEN amount WHEN type = 1 THEN -amount ELSE 0 END) as current
FROM transactions
GROUP BY user_id;