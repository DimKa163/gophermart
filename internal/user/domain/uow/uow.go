package uow

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
)

type UnitOfWork interface {
	UserRepository() repository.UserRepository
	Begin(ctx context.Context) (TxUnitOfWork, error)
}

type TxUnitOfWork interface {
	UserRepository() repository.UserRepository

	Commit(ctx context.Context) error

	Rollback(ctx context.Context) error
}
