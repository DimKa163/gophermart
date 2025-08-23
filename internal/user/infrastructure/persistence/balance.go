package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
)

const (
	balanceGetSQL = `SELECT user_id, current, accrued, withdrawn FROM bonus_balances WHERE user_id = $1`
)

type bonusBalanceRepository struct {
	db db.QueryExecutor
	*db.RetryStrategy
}

func (b *bonusBalanceRepository) Get(ctx context.Context, userID int64) (*model.BonusBalance, error) {
	var balance model.BonusBalance
	var err error
	var currentStr string
	var accrued string
	var withdrawnStr string
	if err = b.QueryRowWithRetry(ctx, b.db, balanceGetSQL, []any{userID}, &balance.UserID, &currentStr, &accrued, &withdrawnStr); err != nil {
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

func NewBonusBalanceRepository(db db.QueryExecutor, retryStrategy *db.RetryStrategy) repository.BonusBalanceRepository {
	return &bonusBalanceRepository{
		db:            db,
		RetryStrategy: retryStrategy,
	}
}
