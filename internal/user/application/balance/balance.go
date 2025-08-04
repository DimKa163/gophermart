package balance

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

type BalanceQuery struct {
}

type BalanceQueryHandler struct {
	uow uow.UnitOfWork
}

func NewBalanceQueryHandler(uow uow.UnitOfWork) *BalanceQueryHandler {
	return &BalanceQueryHandler{uow: uow}
}

func (b *BalanceQueryHandler) Handle(ctx context.Context, _ *BalanceQuery) (*types.AppResult[*model.BonusBalance], error) {
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	bal, err := b.uow.BonusBalanceRepository().Get(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &types.AppResult[*model.BonusBalance]{Code: types.NoChange, Payload: &model.BonusBalance{}}, nil
		}
		return nil, err
	}
	return &types.AppResult[*model.BonusBalance]{Code: types.NoChange, Payload: bal}, nil
}
