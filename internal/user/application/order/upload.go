package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

var ErrOrderConflict = errors.New("order conflict")

type UploadOrderCommand struct {
	Id int64
}

type UploadOrderHandler struct {
	unitOfWork uow.UnitOfWork
}

func NewUploadOrderHandler(unitOfWork uow.UnitOfWork) *UploadOrderHandler {
	return &UploadOrderHandler{unitOfWork}
}

func (handler *UploadOrderHandler) Handle(ctx context.Context, command *UploadOrderCommand) (*types.AppResult[any], error) {
	orderID, err := model.NewOrderID(command.Id)
	if err != nil {
		return &types.AppResult[any]{
			Code:  types.Problem,
			Error: err,
		}, nil
	}

	userID, ok := ctx.Value("userId").(int64)
	if !ok {
		return nil, fmt.Errorf("userId not found in context")
	}

	ord, err := handler.unitOfWork.OrderRepository().Get(ctx, orderID)

	ordExists := true

	if err != nil {
		ordExists = false
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, err
		}
	}

	if ordExists {
		return handler.handleExisted(ctx, userID, ord)
	} else {
		return handler.handleCreated(ctx, userID, orderID)
	}
}

func (handler *UploadOrderHandler) handleExisted(ctx context.Context, userID int64, order *model.Order) (*types.AppResult[any], error) {
	if order.UserID != userID {
		return &types.AppResult[any]{
			Code:  types.Duplicate,
			Error: ErrOrderConflict,
		}, nil
	}
	return &types.AppResult[any]{Code: types.NoChange}, nil
}

func (handler *UploadOrderHandler) handleCreated(ctx context.Context, userID int64, orderID model.OrderID) (*types.AppResult[any], error) {
	txUow, err := handler.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}
	rep := txUow.OrderRepository()
	_, err = rep.Insert(ctx, &model.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  model.OrderStatusNEW,
	})
	if err != nil {
		return nil, err
	}
	if err = txUow.Commit(ctx); err != nil {
		return nil, err
	}
	return &types.AppResult[any]{Code: types.Created}, nil
}
