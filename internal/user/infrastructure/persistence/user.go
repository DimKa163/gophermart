package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
)

const (
	userBalanceSQL = "SELECT user_id, current, accrued, withdrawn FROM bonus_balances WHERE user_id = $1"
	userGetSQL     = "SELECT id, created_at, login, password, salt FROM users WHERE login = $1"
	insertUserSQL  = `INSERT INTO users (created_at, login, password, salt) VALUES ($1, $2, $3, $4) RETURNING id`
	userCountSQL   = `SELECT COUNT(id) FROM users WHERE login = $1`
)

type userRepository struct {
	db db.QueryExecutor
	*db.RetryStrategy
}

func (u *userRepository) GetBonusBalanceByUserID(ctx context.Context, userID int64) (*model.BonusBalance, error) {
	var balance model.BonusBalance
	var err error
	var currentStr string
	var accrued string
	var withdrawnStr string
	if err = u.QueryRowWithRetry(ctx, u.db, userBalanceSQL, []any{userID},
		&balance.UserID,
		&currentStr,
		&accrued,
		&withdrawnStr); err != nil {
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
	var entity model.User
	if err := u.QueryRowWithRetry(ctx, u.db, userGetSQL, []any{login},
		&entity.ID,
		&entity.CreatedAt,
		&entity.Login,
		&entity.Password,
		&entity.Salt); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (u *userRepository) Insert(ctx context.Context, user *model.User) (int64, error) {
	var id int64
	if err := u.QueryRowWithRetry(ctx, u.db, insertUserSQL, []any{
		user.CreatedAt,
		user.Login,
		user.Password,
		user.Salt,
	}, &id); err != nil {
		return -1, err
	}
	return id, nil
}

func (u *userRepository) LoginExists(ctx context.Context, login string) (bool, error) {
	var count int64
	if err := u.db.QueryRow(ctx, userCountSQL, login).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func NewUserRepository(db db.QueryExecutor, retryStrategy *db.RetryStrategy) repository.UserRepository {
	return &userRepository{
		db:            db,
		RetryStrategy: retryStrategy,
	}
}
