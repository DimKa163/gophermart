package balance

import (
	"context"
	"fmt"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

type BalanceQuery struct {
}

type BalanceQueryHandler struct {
	uow uow.UnitOfWork
}

func NewBalanceQueryHandler(uow uow.UnitOfWork) *BalanceQueryHandler {
	return &BalanceQueryHandler{uow: uow}
}

func (b *BalanceQueryHandler) Handle(ctx context.Context, query *BalanceQuery) (*types.AppResult[*model.BonusBalance], error) {
	userID, ok := ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("userId not found in context")
	}
	bal, err := b.uow.BonusBalanceRepository().Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &types.AppResult[*model.BonusBalance]{Code: types.NoChange, Payload: bal}, nil
}
