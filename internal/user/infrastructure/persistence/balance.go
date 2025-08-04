package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
)

type bonusBalanceRepository struct {
	db db.QueryExecutor
}

func (b *bonusBalanceRepository) Get(ctx context.Context, userID int64) (*model.BonusBalance, error) {
	sql := "SELECT user_id, current, accrued, withdrawn FROM bonus_balances WHERE user_id = $1"
	var balance model.BonusBalance
	var err error
	var currentStr string
	var accrued string
	var withdrawnStr string
	if err = b.db.QueryRow(ctx, sql, userID).Scan(&balance.UserID, &currentStr, &accrued, &withdrawnStr); err != nil {
		return nil, err
	}
	balance.Current, err = types.NewDecimalFromString(currentStr)
	if err != nil {
		return nil, err
	}
	balance.Accrued, err = types.NewDecimalFromString(accrued)
	if err != nil {
		return nil, err
	}
	balance.Withdrawn, err = types.NewDecimalFromString(withdrawnStr)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func NewBonusBalanceRepository(db db.QueryExecutor) repository.BonusBalanceRepository {
	return &bonusBalanceRepository{
		db: db,
	}
}
