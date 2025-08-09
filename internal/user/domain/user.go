package domain

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type UserService interface {
	Register(ctx context.Context, login string, password string) (string, error)

	Login(ctx context.Context, login string, password string) (string, error)

	Balance(ctx context.Context) (*model.BonusBalance, error)

	Withdrawal(ctx context.Context) ([]*model.Transaction, error)
}
