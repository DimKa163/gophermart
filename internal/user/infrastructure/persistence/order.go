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

func (o *orderRepository) Exists(ctx context.Context, id model.OrderID) (bool, error) {
	sql := "SELECT COUNT(*) FROM orders WHERE id = $1"
	var count int
	if err := o.db.QueryRow(ctx, sql, id.Value).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (o *orderRepository) Update(ctx context.Context, order *model.Order) error {
	sql := "UPDATE orders SET status=$1, accrual=$2 WHERE id=$3"
	if _, err := o.db.Exec(ctx, sql, order.Status, &order.Accrual, order.OrderID.Value); err != nil {
		return err
	}
	trSQL := "INSERT INTO transactions (created_at, user_id, type, amount, order_id) VALUES ($1, $2, $3, $4, $5)"
	for _, tr := range order.Transactions() {
		if _, err := o.db.Exec(ctx, trSQL, time.Now(), tr.UserID, tr.Type, tr.Amount, tr.OrderID.Value); err != nil {
			return err
		}
	}
	return nil
}

func (o *orderRepository) GetForUpdate(ctx context.Context, limit, offset int, status ...model.OrderStatus) ([]*model.Order, error) {
	sql := "SELECT id, uploaded_at, user_id, status, accrual " +
		"FROM orders WHERE status=ANY($1) ORDER BY uploaded_at LIMIT $2 OFFSET $3 FOR UPDATE SKIP LOCKED"
	var orders []*model.Order
	rows, err := o.db.Query(ctx, sql, status, limit, offset)
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

		order.Accrual = acc
		orders = append(orders, &order)
	}
	return orders, nil
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
	order.Accrual = acc
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

		order.Accrual = acc
		orders = append(orders, &order)
	}
	return orders, nil
}

func (o *orderRepository) Insert(ctx context.Context, order *model.Order) (model.OrderID, error) {
	sql := "INSERT INTO orders (id, uploaded_at, user_id, status) VALUES ($1, $2, $3, $4) RETURNING id"
	var id string
	if err := o.db.QueryRow(ctx, sql, order.OrderID.Value, time.Now(), order.UserID, order.Status).Scan(&id); err != nil {
		return model.DefaultOrderID, err
	}
	orderID, _ := model.NewOrderID(id)
	return orderID, nil
}

func NewOrderRepository(db db.QueryExecutor) repository.OrderRepository {
	return &orderRepository{db}
}
