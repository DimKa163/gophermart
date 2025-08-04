package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type BonusMovementRepository interface {
	GetAll(ctx context.Context, userID int64, tt *model.TransactionType) ([]*model.Transaction, error)
}
