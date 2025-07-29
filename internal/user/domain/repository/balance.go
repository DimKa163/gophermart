package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type BonusBalanceRepository interface {
	Get(ctx context.Context, userId int64) (*model.BonusBalance, error)

	GetForUpdate(ctx context.Context, userId int64) (*model.BonusBalance, error)

	Insert(ctx context.Context, bonus *model.BonusBalance) error

	Update(ctx context.Context, bonus *model.BonusBalance) error
}
