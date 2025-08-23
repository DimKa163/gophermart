package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

type unitOfWork struct {
	db            db.QueryExecutor
	retryStrategy *db.RetryStrategy
	attempts      []int
}

func (u *unitOfWork) BeginTx(ctx context.Context, fn func(ctx context.Context, uow uow.UnitOfWork) error) error {
	err := u.retryStrategy.BeginTx(ctx, u.db, func(ctx context.Context, tx pgx.Tx) error {
		err := fn(ctx, NewUnitOfWork(tx, u.retryStrategy))
		if err != nil {
			_ = tx.Rollback(ctx)
			return err
		}
		if err = tx.Commit(ctx); err != nil {
			return err
		}
		return nil
	})
	return err
}
func (u *unitOfWork) BonusBalanceRepository() repository.BonusBalanceRepository {
	return NewBonusBalanceRepository(u.db, u.retryStrategy)
}
func (u *unitOfWork) BonusMovementRepository() repository.TransactionRepository {
	return NewBonusMovementRepository(u.db, u.retryStrategy)
}
func (u *unitOfWork) UserRepository() repository.UserRepository {
	return NewUserRepository(u.db, u.retryStrategy)
}
func (u *unitOfWork) OrderRepository() repository.OrderRepository {
	return NewOrderRepository(u.db, u.retryStrategy)
}
func NewUnitOfWork(db db.QueryExecutor, retryStrategy *db.RetryStrategy) uow.UnitOfWork {
	return &unitOfWork{
		db:            db,
		retryStrategy: retryStrategy,
	}
}
