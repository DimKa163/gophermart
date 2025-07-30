package withdrawal

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

type WithdrawalQuery struct{}

type WithdrawalQueryHandler struct {
	uow uow.UnitOfWork
}

func NewWithdrawalQueryHandler(uow uow.UnitOfWork) *WithdrawalQueryHandler {
	return &WithdrawalQueryHandler{uow: uow}
}

func (w *WithdrawalQueryHandler) Handle(ctx context.Context, _ *WithdrawalQuery) (*types.AppResult[[]*model.BonusMovement], error) {
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	rep := w.uow.BonusMovementRepository()
	t := model.WITHDRAWAL
	items, err := rep.GetAll(ctx, userID, &t)
	if err != nil {
		return nil, err
	}
	return &types.AppResult[[]*model.BonusMovement]{
		Code:    types.NoChange,
		Payload: items,
	}, nil
}
