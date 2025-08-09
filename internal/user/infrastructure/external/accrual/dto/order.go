package dto

import (
	"github.com/DimKa163/gophermart/internal/shared/types"
)

type OrderStatus string

const (
	StatusREGISTERED OrderStatus = "REGISTERED"
	StatusPROCESSING OrderStatus = "PROCESSING"
	StatusPROCESSED  OrderStatus = "PROCESSED"
	StatusINVALID    OrderStatus = "INVALID"
)

func (s OrderStatus) String() string {
	return string(s)
}

type Order struct {
	Number  string
	Status  OrderStatus
	Accrual *types.Decimal
}
