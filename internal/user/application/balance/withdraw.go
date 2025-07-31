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

var ErrNegativeBalance = errors.New("not enough balance")

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
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	txUow, err := wh.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	balRep := txUow.BonusBalanceRepository()
	defer func() {
		if err != nil {
			_ = txUow.Rollback(ctx)
			return
		}
		_ = txUow.Commit(ctx)
		return
	}()
	balance, err := balRep.GetForUpdate(ctx, userID)
	if err != nil {
		return nil, err
	}

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

	if balance.Current.Cmp(command.Sum) == -1 {
		return &types.AppResult[any]{
			Code:  types.Problem,
			Error: ErrNegativeBalance,
		}, nil
	}
	movement, err := model.NewBonusMovement(userID, model.WITHDRAWAL, command.Sum, order.OrderID)
	if err != nil {
		return nil, err
	}
	movRep := txUow.BonusMovementRepository()
	if err = movRep.Insert(ctx, movement); err != nil {
		return nil, err
	}
	accSum, err := movRep.Sum(ctx, userID, model.ACCRUAL)
	if err != nil {
		return nil, err
	}
	withSum, err := movRep.Sum(ctx, userID, model.WITHDRAWAL)
	if err != nil {
		return nil, err
	}
	result := accSum.Sub(*withSum)
	balance.Current = result
	balance.Withdrawn = *withSum
	if err = balRep.Update(ctx, balance); err != nil {
		return nil, err
	}
	return &types.AppResult[any]{
		Code: types.Created,
	}, nil
}
