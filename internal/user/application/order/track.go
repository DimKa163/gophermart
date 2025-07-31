package order

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

type TrackOrderCommand struct {
	Limit int
}

type TrackOrderHandler struct {
	uow uow.UnitOfWork
}

func NewTrackOrderHandler(uow uow.UnitOfWork) *TrackOrderHandler {
	return &TrackOrderHandler{uow: uow}
}

func (handler *TrackOrderHandler) Handle(ctx context.Context, command *TrackOrderCommand) (*types.AppResult[any], error) {
	var err error
	txUow, err := handler.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = txUow.Rollback(ctx)
			return
		}
		_ = txUow.Commit(ctx)
	}()
	orderRep := txUow.OrderRepository()
	_, err = orderRep.GetForUpdate(ctx, command.Limit, model.OrderStatusNEW, model.OrderStatusPROCESSING)
	if err != nil {
		return nil, err
	}
	//TODO call accrual system

	return &types.AppResult[any]{}, nil
}
