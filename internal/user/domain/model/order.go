package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"strconv"
	"time"
	"unicode"
)

type OrderStatus int

const (
	OrderStatusNEW OrderStatus = iota
	OrderStatusPROCESSING
	OrderStatusINVALID
	OrderStatusPROCESSED
)

func (s OrderStatus) String() string {
	return [...]string{"NEW", "PROCESSING", "INVALID", "PROCESSED"}[s]
}

type Order struct {
	OrderID    OrderID
	UploadedAt *time.Time
	Status     OrderStatus
	UserID     int64
	Accrual    *types.Decimal
}

var ErrOrderID = errors.New("order code not valid")

var DefaultOrderID = OrderID{
	Value: 0,
}

type OrderID struct {
	Value int64
}

func NewOrderID(value int64) (OrderID, error) {
	if err := validate(value); err != nil {
		return DefaultOrderID, err
	}
	return OrderID{value}, nil
}

func (id *OrderID) MarshalJSON() ([]byte, error) {
	return []byte(id.String()), nil
}

func (id *OrderID) UnmarshalJSON(data []byte) error {
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return fmt.Errorf("OrderID.UnmarshalJSON: %w", err)
	}
	val, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return fmt.Errorf("OrderID.UnmarshalJSON: %w", err)
	}
	err = validate(val)
	if err != nil {
		return fmt.Errorf("OrderID.UnmarshalJSON: %w", err)
	}
	id.Value = val
	return nil
}
func validate(value int64) error {
	v := strconv.FormatInt(value, 10)
	num := v
	var sum int
	var double bool
	for i := len(num) - 1; i >= 0; i-- {
		r := rune(num[i])
		if !unicode.IsNumber(r) {
			return ErrOrderID
		}
		d := int(r - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	if sum%10 != 0 {
		return ErrOrderID
	}
	return nil
}

func (id *OrderID) String() string {
	return strconv.FormatInt(id.Value, 10)
}
