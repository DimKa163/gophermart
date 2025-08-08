package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type BonusBalanceRepository interface {
	Get(ctx context.Context, userID int64) (*model.BonusBalance, error)
}
