package uow

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
)

type UnitOfWork interface {
	UserRepository() repository.UserRepository
	OrderRepository() repository.OrderRepository
	BonusBalanceRepository() repository.BonusBalanceRepository
	BonusMovementRepository() repository.TransactionRepository

	BeginTx(ctx context.Context, fn func(ctx context.Context, uow UnitOfWork) error) error
}
