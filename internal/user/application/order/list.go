package order

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

type OrderQuery struct {
}
type OrderQueryHandler struct {
	uow uow.UnitOfWork
}

func NewOrderQueryHandler(uow uow.UnitOfWork) *OrderQueryHandler {
	return &OrderQueryHandler{uow: uow}
}

func (h *OrderQueryHandler) Handle(ctx context.Context, _ *OrderQuery) (*types.AppResult[[]*model.Order], error) {
	rep := h.uow.OrderRepository()
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	orders, err := rep.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &types.AppResult[[]*model.Order]{
		Code:    types.NoChange,
		Payload: orders,
	}, nil
}
