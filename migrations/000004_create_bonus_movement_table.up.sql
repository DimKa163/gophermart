CREATE TABLE IF NOT EXISTS bonus_movements
(
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    user_id BIGSERIAL REFERENCES users(id),
    type INT,
    amount DECIMAL(10, 2),
    order_id BIGSERIAL REFERENCES orders(id)
);

CREATE INDEX IF NOT EXISTS bonus_movement_user_id_ix ON bonus_movements(user_id ASC);