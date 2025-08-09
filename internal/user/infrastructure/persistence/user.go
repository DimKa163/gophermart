package persistence

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db db.QueryExecutor
}

func (u *userRepository) GetBonusBalanceByUserID(ctx context.Context, userID int64) (*model.BonusBalance, error) {
	sql := "SELECT user_id, current, accrued, withdrawn FROM bonus_balances WHERE user_id = $1"
	var balance model.BonusBalance
	var err error
	var currentStr string
	var accrued string
	var withdrawnStr string
	if err = u.db.QueryRow(ctx, sql, userID).Scan(&balance.UserID, &currentStr, &accrued, &withdrawnStr); err != nil {
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
func (u *userRepository) Get(ctx context.Context, login string) (*model.User, error) {
	sql := "SELECT id, created_at, login, password, salt FROM users WHERE login = $1"
	var entity model.User
	if err := u.db.QueryRow(ctx, sql, login).Scan(
		&entity.ID,
		&entity.CreatedAt,
		&entity.Login,
		&entity.Password,
		&entity.Salt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return &entity, nil
}

func (u *userRepository) Insert(ctx context.Context, user *model.User) (int64, error) {
	sql := "INSERT INTO users (created_at, login, password, salt) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int64
	if err := u.db.QueryRow(
		ctx,
		sql,
		user.CreatedAt,
		user.Login,
		user.Password,
		user.Salt,
	).Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (u *userRepository) LoginExists(ctx context.Context, login string) (bool, error) {
	var count int64
	sql := "SELECT COUNT(id) FROM users WHERE login = $1"
	if err := u.db.QueryRow(ctx, sql, login).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func NewUserRepository(db db.QueryExecutor) repository.UserRepository {
	return &userRepository{db}
}
