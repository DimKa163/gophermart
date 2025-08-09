DROP VIEW IF EXISTS bonus_balances;

DROP INDEX IF EXISTS transactions_user_id_ix ON public.transactions;

DROP TABLE IF EXISTS transactions;