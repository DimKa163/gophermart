package model

import (
	"database/sql/driver"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"time"
)

type BonusMovementType int

var ErrBonusMovement = errors.New("invalid bonus movement")

const (
	ACCRUAL BonusMovementType = iota
	WITHDRAWAL
)

func (s *BonusMovementType) String() string {
	return [...]string{"ACCRUAL", "WITHDRAWAL"}[*s]
}

func (s *BonusMovementType) Scan(value interface{}) error {
	switch value.(type) {
	case int:
		*s = BonusMovementType(value.(int))
		break
	}
	return nil
}

func (s *BonusMovementType) Value() (driver.Value, error) {
	return int64(*s), nil
}

type BonusMovement struct {
	UserID    int64
	CreatedAt time.Time
	Type      BonusMovementType
	Amount    types.Decimal
	OrderID   OrderID
}

func NewBonusMovement(userID int64, tt BonusMovementType, amount types.Decimal, orderId OrderID) (*BonusMovement, error) {
	if amount.IsNegative() {
		return nil, ErrBonusMovement
	}
	return &BonusMovement{
		UserID:    userID,
		CreatedAt: time.Now(),
		Type:      tt,
		Amount:    amount,
		OrderID:   orderId,
	}, nil
}
