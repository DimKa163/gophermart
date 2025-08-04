package balance

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

var ErrNegativeBalance = errors.New("not enough bal")

var ErrWrongOrder = errors.New("wrong order")

type WithdrawCommand struct {
	OrderID model.OrderID `json:"order_id"`
	Sum     types.Decimal `json:"amount"`
}

type WithdrawHandler struct {
	uow uow.UnitOfWork
}

func NewWithdrawHandler(uow uow.UnitOfWork) *WithdrawHandler {
	return &WithdrawHandler{uow: uow}
}

func (wh *WithdrawHandler) Handle(ctx context.Context, command *WithdrawCommand) (*types.AppResult[any], error) {
	var err error
	txUow, err := wh.uow.Begin(ctx)
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

	order, err := orderRep.Get(ctx, command.OrderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &types.AppResult[any]{
				Code:  types.Problem,
				Error: ErrWrongOrder,
			}, nil
		}
		return nil, err
	}
	bal, err := txUow.UserRepository().GetBonusBalanceByUserID(ctx, order.UserID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if bal == nil {
		bal = &model.BonusBalance{}
	}
	if bal.Current.Cmp(command.Sum) < 0 {
		return nil, ErrNegativeBalance
	}
	order.AddTransaction(model.WITHDRAWAL, command.Sum)
	err = orderRep.Update(ctx, order)
	if err != nil {
		return nil, err
	}
	return &types.AppResult[any]{
		Code: types.Created,
	}, nil
}
