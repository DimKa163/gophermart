package order

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

type OrderQuery struct {
	UserID int64
}
type OrderQueryHandler struct {
	uow uow.UnitOfWork
}

func NewOrderQueryHandler(uow uow.UnitOfWork) *OrderQueryHandler {
	return &OrderQueryHandler{uow: uow}
}

func (h *OrderQueryHandler) Handle(ctx context.Context, query OrderQuery) (*types.AppResult[[]*model.Order], error) {
	rep := h.uow.OrderRepository()
	orders, err := rep.GetAll(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	return &types.AppResult[[]*model.Order]{
		Code:    types.NoChange,
		Payload: orders,
	}, nil
}
