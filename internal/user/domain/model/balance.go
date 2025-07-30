package model

import (
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"time"
)

var ErrBonusBalance = errors.New("invalid bonus balance")

type BonusBalance struct {
	UserID    int64
	CreatedAt time.Time
	Current   types.Decimal
	Withdrawn types.Decimal
}

func NewBonusBalance(userID int64, current, withdrawn types.Decimal) (*BonusBalance, error) {
	if current.IsNegative() || withdrawn.IsNegative() {
		return nil, ErrBonusBalance
	}
	return &BonusBalance{
		UserID:    userID,
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}
