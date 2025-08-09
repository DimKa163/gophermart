package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type OrderRepository interface {
	Exists(ctx context.Context, id model.OrderID) (bool, error)

	Get(ctx context.Context, id model.OrderID) (*model.Order, error)

	GetForUpdate(ctx context.Context, limit, offset int, status ...model.OrderStatus) ([]*model.Order, error)

	GetAll(ctx context.Context, userID int64) ([]*model.Order, error)

	Insert(ctx context.Context, order *model.Order) (model.OrderID, error)

	Update(ctx context.Context, order *model.Order) error
}
