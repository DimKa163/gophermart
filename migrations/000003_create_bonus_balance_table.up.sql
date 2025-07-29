CREATE TABLE IF NOT EXISTS bonus_balances
(
    user_id bigserial PRIMARY KEY REFERENCES users(id),
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    current DECIMAL(10, 2),
    withdrawn DECIMAL(10, 2)
)