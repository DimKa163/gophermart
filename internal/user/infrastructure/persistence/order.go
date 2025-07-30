package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"time"
)

type orderRepository struct {
	db db.QueryExecutor
}

func (o *orderRepository) Get(ctx context.Context, id model.OrderID) (*model.Order, error) {
	sql := "SELECT id, uploaded_at, user_id, status, accrual FROM orders WHERE id=$1"
	var order model.Order
	var orderID int64
	var accrual *string
	if err := o.db.QueryRow(ctx, sql, id.Value).Scan(&orderID, &order.UploadedAt, &order.UserID, &order.Status, &accrual); err != nil {
		return nil, err
	}
	var acc types.Decimal
	var err error
	if accrual != nil {
		acc, err = types.NewDecimalFromString(*accrual)
		if err != nil {
			return nil, err
		}
	}
	order.Accrual = &acc
	order.OrderID = model.OrderID{Value: orderID}
	return &order, nil
}

func (o *orderRepository) GetAll(ctx context.Context, userID int64) ([]*model.Order, error) {
	sql := "SELECT id, uploaded_at, user_id, status, accrual FROM orders WHERE user_id=$1"
	var orders []*model.Order
	rows, err := o.db.Query(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderID int64
		var order model.Order
		var accrual *string
		if err := rows.Scan(&orderID, &order.UploadedAt, &order.UserID, &order.Status, &accrual); err != nil {
			return nil, err
		}
		order.OrderID = model.OrderID{Value: orderID}
		var acc types.Decimal
		if accrual != nil {
			acc, err = types.NewDecimalFromString(*accrual)
			if err != nil {
				return nil, err
			}
		}

		order.Accrual = &acc
		orders = append(orders, &order)
	}
	return orders, nil
}

func (o *orderRepository) Insert(ctx context.Context, order *model.Order) (string, error) {
	sql := "INSERT INTO orders (id, uploaded_at, user_id, status) VALUES ($1, $2, $3, $4) RETURNING id"
	var id string
	if err := o.db.QueryRow(ctx, sql, order.OrderID.Value, time.Now(), order.UserID, order.Status).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

func NewOrderRepository(db db.QueryExecutor) repository.OrderRepository {
	return &orderRepository{db}
}
