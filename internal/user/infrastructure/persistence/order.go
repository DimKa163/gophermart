package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type orderRepository struct {
	db db.QueryExecutor
}

func (o *orderRepository) Get(ctx context.Context, id model.OrderID) (*model.Order, error) {
	sql := "SELECT id, uploaded_at, user_id, status, accrual FROM orders WHERE id=$1"
	var order model.Order
	var orderId int64
	var accrual pgtype.Numeric
	if err := o.db.QueryRow(ctx, sql, id.Value).Scan(&orderId, &order.UploadedAt, &order.UserID, &order.Status, &accrual); err != nil {
		return nil, err
	}
	if err := accrual.Scan(accrual); err != nil {
		return nil, err
	}
	order.OrderID = model.OrderID{Value: orderId}
	return &order, nil
}

func (o *orderRepository) GetAll(ctx context.Context, userId int64) ([]*model.Order, error) {
	sql := "SELECT id, uploaded_at, user_id, status, accrual FROM orders WHERE user_id=$1"
	var orders []*model.Order
	rows, err := o.db.Query(ctx, sql, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var orderId int64
		var order model.Order
		var accrual pgtype.Numeric
		if err := rows.Scan(&orderId, &order.UploadedAt, &order.UserID, &order.Status, &accrual); err != nil {
			return nil, err
		}
		if err := accrual.Scan(order.Accrual); err != nil {
			return nil, err
		}
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
