DROP INDEX IF EXISTS orders_status_ix ON public.orders;
DROP INDEX IF EXISTS orders_user_id_ix ON public.orders;

DROP TABLE IF EXISTS orders;