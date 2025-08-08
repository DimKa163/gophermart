package application

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

var (
	ErrOrderIdProblem             = errors.New("Order ID is invalid")
	ErrOrderExistsWithAnotherUser = &domain.ResourceAlreadyExists{Message: "Order already exists"}
	ErrNegativeBalance            = domain.NewProblemError("Not enough bonus points", nil)
)

type orderService struct {
	uow uow.UnitOfWork
}

func (o *orderService) Upload(ctx context.Context, number string) (bool, error) {
	orderID, err := model.NewOrderID(number)
	if err != nil {
		return false, err
	}
	userID, err := auth.User(ctx)
	if err != nil {
		return false, err
	}
	orderRep := o.uow.OrderRepository()
	ord, err := orderRep.Get(ctx, orderID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}
	if ord != nil {
		if ord.UserID != userID {
			return false, ErrOrderExistsWithAnotherUser
		}
		return false, nil
	}
	_, err = orderRep.Insert(ctx, &model.Order{
		OrderID: orderID,
		UserID:  userID,
		Status:  model.OrderStatusNEW,
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

func (o *orderService) List(ctx context.Context) ([]*model.Order, error) {
	rep := o.uow.OrderRepository()
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	orders, err := rep.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (o *orderService) Withdraw(ctx context.Context, number string, sum types.Decimal) error {
	userID, err := auth.User(ctx)
	if err != nil {
		return err
	}
	orderID, err := model.NewOrderID(number)
	if err != nil {
		return err
	}
	orderRep := o.uow.OrderRepository()
	ord, err := orderRep.Get(ctx, orderID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		if _, err = o.Upload(ctx, number); err != nil {
			return err
		}
		if ord, err = orderRep.Get(ctx, orderID); err != nil {
			return err
		}
	}
	userRep := o.uow.UserRepository()
	bal, err := userRep.GetBonusBalanceByUserID(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if bal == nil {
		bal = &model.BonusBalance{}
	}
	if bal.Current.Cmp(sum) < 0 {
		return ErrNegativeBalance
	}
	ord.AddTransaction(model.WITHDRAWAL, sum)
	if err = orderRep.Update(ctx, ord); err != nil {
		return err
	}
	return nil
}

func NewOrderService(uow uow.UnitOfWork) domain.OrderService {
	return &orderService{uow: uow}
}
