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

var ErrNegativeBalance = errors.New("not enough bal")

var ErrWrongOrder = errors.New("wrong order")

type WithdrawCommand struct {
	OrderID string        `json:"order_id"`
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
	orderID, err := model.NewOrderID(command.OrderID)
	if err != nil {
		return &types.AppResult[any]{
			Code:  types.Problem,
			Error: ErrWrongOrder,
		}, nil
	}
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
	order, err := wh.getOrCreate(ctx, txUow, userID, orderID)
	if err != nil {
		return nil, err
	}
	orderRep := txUow.OrderRepository()

	bal, err := txUow.UserRepository().GetBonusBalanceByUserID(ctx, userID)
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

func (wh *WithdrawHandler) getOrCreate(ctx context.Context, txUow uow.TxUnitOfWork, userID int64, orderID model.OrderID) (*model.Order, error) {
	var err error
	orderRep := txUow.OrderRepository()
	ex, err := orderRep.Exists(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if !ex {
		orderID, err = orderRep.Insert(ctx, &model.Order{OrderID: orderID,
			UserID: userID,
			Status: model.OrderStatusNEW})
		if err != nil {
			return nil, err
		}
	}
	return orderRep.Get(ctx, orderID)
}
