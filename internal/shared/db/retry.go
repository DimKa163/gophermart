package db

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v5"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type RetryStrategy struct {
	attempts []int
}

func NewRetryStrategy(attempts []int) *RetryStrategy {
	return &RetryStrategy{attempts}
}

func (s *RetryStrategy) BeginTx(ctx context.Context, db QueryExecutor, fn func(ctx context.Context, tx pgx.Tx) error) error {
	seconds := s.attempts
	attempt := 0
	_, err := backoff.Retry(ctx, func() (bool, error) {
		tx, err := db.Begin(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if ShouldRetry(pgErr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return false, backoff.RetryAfter(seconds[at])
				}
			}
			return false, backoff.Permanent(err)
		}
		err = fn(ctx, tx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if ShouldRetry(pgErr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return false, backoff.RetryAfter(seconds[at])
				}
			}
			return false, backoff.Permanent(err)
		}
		return true, nil
	})
	return err
}

func (s *RetryStrategy) ExecWithRetry(ctx context.Context, fn func(ctx context.Context) (pgconn.CommandTag, error)) (pgconn.CommandTag, error) {
	seconds := s.attempts
	attempt := 0
	return backoff.Retry(ctx, func() (pgconn.CommandTag, error) {
		tag, err := fn(ctx)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if ShouldRetry(pgErr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return pgconn.CommandTag{}, backoff.RetryAfter(seconds[at])
				}
			}
			return pgconn.CommandTag{}, backoff.Permanent(err)
		}
		return tag, nil
	})
}

func (s *RetryStrategy) QueryWithRetry(ctx context.Context, db QueryExecutor, sql string, args ...any) (pgx.Rows, error) {
	seconds := s.attempts
	attempt := 0
	return backoff.Retry(ctx, func() (pgx.Rows, error) {
		rows, err := db.Query(ctx, sql, args...)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if ShouldRetry(pgErr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return nil, backoff.RetryAfter(seconds[at])
				}
			}
			return nil, backoff.Permanent(err)
		}
		return rows, nil
	})
}
func (s *RetryStrategy) QueryRowWithRetry(ctx context.Context, db QueryExecutor, sql string, args []any, dest ...any) error {
	seconds := s.attempts
	attempt := 0
	if _, err := backoff.Retry(ctx, func() (bool, error) {
		row := db.QueryRow(ctx, sql, args...)
		if err := row.Scan(dest...); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if ShouldRetry(pgErr) && attempt < len(seconds) {
					at := attempt
					attempt++
					return false, backoff.RetryAfter(seconds[at])
				}
			}
			return false, backoff.Permanent(err)
		}
		return true, nil
	}); err != nil {
		return err
	}
	return nil
}

func ShouldRetry(pgErr *pgconn.PgError) bool {
	switch pgErr.Code {
	case pgerrcode.SerializationFailure:
	case pgerrcode.DeadlockDetected:
	case pgerrcode.TooManyConnections:
	case pgerrcode.LockNotAvailable:
	case pgerrcode.CannotConnectNow:
	case pgerrcode.QueryCanceled:
	case pgerrcode.UniqueViolation:
		return true
	}
	return false
}
