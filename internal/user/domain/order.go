package domain

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type OrderService interface {
	Upload(ctx context.Context, number string) (bool, error)

	List(ctx context.Context) ([]*model.Order, error)

	Withdraw(ctx context.Context, number string, decimal types.Decimal) error
}
