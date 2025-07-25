package model

import (
	"errors"
	"time"
)

type BonusMovementType int

var ErrBonusMovement = errors.New("invalid bonus movement")

const (
	ACCRUAL BonusMovementType = iota
	WITHDRAWAL
)

func (s BonusMovementType) String() string {
	return [...]string{"ACCRUAL", "WITHDRAWAL"}[s]
}

type BonusMovement struct {
	UserID    int64
	CreatedAt time.Time
	Type      BonusMovementType
	Amount    float64
	OrderID   *OrderID
}

func NewBonusMovement(userID int64, tt BonusMovementType, amount float64, orderId *OrderID) (*BonusMovement, error) {
	if amount < 0 {
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
