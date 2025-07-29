package persistence

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/jackc/pgx/v5"
)

type userRepository struct {
	db db.QueryExecutor
}

func (u *userRepository) Get(ctx context.Context, login string) (*model.User, error) {
	sql := "SELECT id, created_at, login, password FROM users WHERE login = $1"
	var entity model.User
	if err := entity.Scan(u.db.QueryRow(ctx, sql, login)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}
	return &entity, nil
}

func (u *userRepository) Insert(ctx context.Context, user *model.User) (int64, error) {
	sql := "INSERT INTO users (created_at, login, password) VALUES ($1, $2, $3) RETURNING id"
	var id int64
	if err := u.db.QueryRow(ctx, sql, user.CreatedAt, user.Login, user.Password).Scan(&id); err != nil {
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
