package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"time"
)

type bonusBalanceRepository struct {
	db db.QueryExecutor
}

func (b *bonusBalanceRepository) GetForUpdate(ctx context.Context, userID int64) (*model.BonusBalance, error) {
	sql := "SELECT user_id, created_at, current, withdrawn FROM bonus_balances WHERE user_id = $1 FOR UPDATE"
	var balance model.BonusBalance
	var err error
	var currentStr string
	var withdrawnStr string
	if err = b.db.QueryRow(ctx, sql, userID).Scan(&balance.UserID, &balance.CreatedAt, &currentStr, &withdrawnStr); err != nil {
		return nil, err
	}
	balance.Current, err = types.NewDecimalFromString(currentStr)
	if err != nil {
		return nil, err
	}
	balance.Withdrawn, err = types.NewDecimalFromString(withdrawnStr)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (b *bonusBalanceRepository) Get(ctx context.Context, userID int64) (*model.BonusBalance, error) {
	sql := "SELECT user_id, created_at, current, withdrawn FROM bonus_balances WHERE user_id = $1"
	var balance model.BonusBalance
	var err error
	var currentStr string
	var withdrawnStr string
	if err = b.db.QueryRow(ctx, sql, userID).Scan(&balance.UserID, &balance.CreatedAt, &currentStr, &withdrawnStr); err != nil {
		return nil, err
	}
	balance.Current, err = types.NewDecimalFromString(currentStr)
	if err != nil {
		return nil, err
	}
	balance.Withdrawn, err = types.NewDecimalFromString(withdrawnStr)
	if err != nil {
		return nil, err
	}
	return &balance, nil
}

func (b *bonusBalanceRepository) Insert(ctx context.Context, bonus *model.BonusBalance) error {
	sql := "INSERT INTO bonus_balances (user_id, created_at, current, withdrawn) VALUES ($1, $2, $3, $4)"
	if _, err := b.db.Exec(ctx, sql, bonus.UserID, time.Now(), bonus.Current, bonus.Withdrawn); err != nil {
		return err
	}
	return nil
}

func (b *bonusBalanceRepository) Update(ctx context.Context, bonus *model.BonusBalance) error {
	sql := "UPDATE bonus_balances SET current = $1, withdrawn = $2 WHERE user_id = $3"
	if _, err := b.db.Exec(ctx, sql, bonus.Current, bonus.Withdrawn, bonus.UserID); err != nil {
		return err
	}
	return nil
}

func NewBonusBalanceRepository(db db.QueryExecutor) repository.BonusBalanceRepository {
	return &bonusBalanceRepository{
		db: db,
	}
}
