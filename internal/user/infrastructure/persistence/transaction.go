package persistence

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/db"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/repository"
	"github.com/jackc/pgx/v5"
)

type bonusMovementRepository struct {
	db db.QueryExecutor
}

func (b bonusMovementRepository) GetAll(ctx context.Context, userID int64, tt *model.TransactionType) ([]*model.Transaction, error) {
	var sql string
	var rows pgx.Rows
	var err error
	if tt != nil {
		sql = "SELECT created_at, user_id, type, amount, order_id FROM transactions WHERE user_id = $1 AND type = $2"
		rows, err = b.db.Query(ctx, sql, userID, *tt)
	} else {
		sql = "SELECT created_at, user_id, type, amount, order_id FROM transactions WHERE user_id = $1"
		rows, err = b.db.Query(ctx, sql, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var transactions []*model.Transaction
	for rows.Next() {
		var transaction model.Transaction
		var orderID int64
		var amountStr string
		if err = rows.Scan(&transaction.CreatedAt, &transaction.UserID, &transaction.Type, &amountStr, &orderID); err != nil {
			return nil, err
		}
		transaction.OrderID = model.OrderID{
			Value: orderID,
		}
		transaction.Amount, err = types.NewDecimalFromString(amountStr)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, &transaction)
	}
	return transactions, nil
}

func NewBonusMovementRepository(db db.QueryExecutor) repository.TransactionRepository {
	return &bonusMovementRepository{
		db: db,
	}
}
