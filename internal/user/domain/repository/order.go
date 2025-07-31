package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type OrderRepository interface {
	Get(ctx context.Context, id model.OrderID) (*model.Order, error)

	GetForUpdate(ctx context.Context, limit int, status ...model.OrderStatus) ([]*model.Order, error)

	GetAll(ctx context.Context, userID int64) ([]*model.Order, error)

	Insert(ctx context.Context, order *model.Order) (string, error)
}
