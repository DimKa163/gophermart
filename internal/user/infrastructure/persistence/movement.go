package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/jackc/pgx/v5"
	"time"
)

type bonusMovementRepository struct {
	db db.QueryExecutor
}

func (b bonusMovementRepository) Sum(ctx context.Context, userId int64, tt model.BonusMovementType) (*types.Decimal, error) {
	sql := "SELECT SUM(amount) FROM bonus_movements WHERE user_id = $1 AND type = $2"
	var amountStr string
	if err := b.db.QueryRow(ctx, sql, userId, tt).Scan(&amountStr); err != nil {
		return nil, err
	}
	amount, err := types.NewDecimalFromString(amountStr)
	if err != nil {
		return nil, err
	}
	return &amount, nil
}

func (b bonusMovementRepository) GetAll(ctx context.Context, userId int64, tt *model.BonusMovementType) ([]*model.BonusMovement, error) {
	var sql string
	var rows pgx.Rows
	var err error
	if tt != nil {
		sql = "SELECT created_at, user_id, type, amount FROM bonus_movements WHERE user_id = $1 AND type = $2"
		rows, err = b.db.Query(ctx, sql, userId, *tt)
	} else {
		sql = "SELECT created_at, user_id, type, amount FROM bonus_movements WHERE user_id = $1"
		rows, err = b.db.Query(ctx, sql, userId)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movements []*model.BonusMovement
	for rows.Next() {
		var movement model.BonusMovement
		var amountStr string
		if err = rows.Scan(&movement.CreatedAt, &movement.UserID, &movement.Type, &amountStr); err != nil {
			return nil, err
		}
		movement.Amount, err = types.NewDecimalFromString(amountStr)
		if err != nil {
			return nil, err
		}
		movements = append(movements, &movement)
	}
	return movements, nil
}

func (b bonusMovementRepository) Insert(ctx context.Context, bonus *model.BonusMovement) error {
	sql := "INSERT INTO bonus_movements (created_at, user_id, type, amount) VALUES ($1, $2, $3, $4)"

	if _, err := b.db.Exec(ctx, sql, time.Now(), bonus.UserID, bonus.Type, bonus.Amount); err != nil {
		return err
	}

	return nil
}

func NewBonusMovementRepository(db db.QueryExecutor) repository.BonusMovementRepository {
	return &bonusMovementRepository{
		db: db,
	}
}
