package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type UserRepository interface {
	Get(ctx context.Context, login string) (*model.User, error)

	LoginExists(ctx context.Context, login string) (bool, error)

	Insert(ctx context.Context, user *model.User) (int64, error)
}
