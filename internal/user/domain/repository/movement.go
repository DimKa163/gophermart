package repository

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
)

type BonusMovementRepository interface {
	GetAll(ctx context.Context, userID int64, tt *model.BonusMovementType) ([]*model.BonusMovement, error)
	Insert(ctx context.Context, bonus *model.BonusMovement) error

	Sum(ctx context.Context, userID int64, tt model.BonusMovementType) (*types.Decimal, error)
}
