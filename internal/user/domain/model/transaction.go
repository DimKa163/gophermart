package model

import (
	"database/sql/driver"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"time"
)

type TransactionType int

var ErrBonusTransaction = errors.New("invalid transaction")

const (
	ACCRUAL TransactionType = iota
	WITHDRAWAL
)

func (s *TransactionType) String() string {
	return [...]string{"ACCRUAL", "WITHDRAWAL"}[*s]
}

func (s *TransactionType) Value() (driver.Value, error) {
	return int64(*s), nil
}

type Transaction struct {
	UserID    int64
	CreatedAt time.Time
	Type      TransactionType
	Amount    types.Decimal
	OrderID   OrderID
}
