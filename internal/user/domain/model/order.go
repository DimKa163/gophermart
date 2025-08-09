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
	OrderID      OrderID
	UploadedAt   *time.Time
	Status       OrderStatus
	UserID       int64
	Withdrawn    types.Decimal
	Accrual      types.Decimal
	transactions []*Transaction
	Error        string
}

func (o *Order) AddTransaction(tt TransactionType, amount types.Decimal) {
	if o.transactions == nil {
		o.transactions = make([]*Transaction, 0)
	}
	o.transactions = append(o.transactions, &Transaction{
		Amount:  amount,
		Type:    tt,
		UserID:  o.UserID,
		OrderID: o.OrderID,
	})
	switch tt {
	case ACCRUAL:
		o.Accrual = o.Accrual.Add(amount)
	case WITHDRAWAL:
		o.Withdrawn = o.Withdrawn.Add(amount)
	}
}

func (o *Order) Transactions() []*Transaction {
	return o.transactions
}

type OrderIDError struct {
	Message string
}

func (e *OrderIDError) Error() string {
	return e.Message
}

var ErrOrderID = errors.New("Order ID is invalid")

var DefaultOrderID = OrderID{
	Value: 0,
}

type OrderID struct {
	Value int64
}

func NewOrderID(value string) (OrderID, error) {
	if err := validate(value); err != nil {
		return DefaultOrderID, err
	}
	v, _ := strconv.ParseInt(value, 10, 64)
	return OrderID{Value: v}, nil
}

func (id OrderID) MarshalJSON() ([]byte, error) {
	return []byte("\"" + id.String() + "\""), nil
}

func (id OrderID) UnmarshalJSON(data []byte) error {
	var strVal string
	if err := json.Unmarshal(data, &strVal); err != nil {
		return fmt.Errorf("OrderID.UnmarshalJSON: %w", err)
	}
	err := validate(strVal)
	if err != nil {
		return fmt.Errorf("OrderID.UnmarshalJSON: %w", err)
	}
	val, err := strconv.ParseInt(strVal, 10, 64)
	if err != nil {
		return err
	}
	id.Value = val
	return nil
}
func (id *OrderID) String() string {
	return strconv.FormatInt(id.Value, 10)
}

func validate(value string) error {
	num := value
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
