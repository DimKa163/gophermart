package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

type unitOfWork struct {
	db db.QueryExecutor
}

func (u *unitOfWork) Begin(ctx context.Context) (uow.TxUnitOfWork, error) {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &txUnitOfWork{
		UnitOfWork: &unitOfWork{
			db: tx,
		},
		db: tx,
	}, nil
}

func (u *unitOfWork) UserRepository() repository.UserRepository {
	return NewUserRepository(u.db)
}

func (u *unitOfWork) OrderRepository() repository.OrderRepository { return NewOrderRepository(u.db) }

func NewUnitOfWork(db db.QueryExecutor) uow.UnitOfWork {
	return &unitOfWork{
		db: db,
	}
}

type txUnitOfWork struct {
	uow.UnitOfWork
	db db.TxQueryExecutor
}

func (u *txUnitOfWork) Commit(ctx context.Context) error {
	return u.db.Commit(ctx)
}

func (u *txUnitOfWork) Rollback(ctx context.Context) error {
	return u.db.Rollback(ctx)
}
