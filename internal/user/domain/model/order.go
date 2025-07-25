package model

import (
	"errors"
	"time"
	"unicode"
)

type Order struct {
	OrderID    OrderID
	UploadedAt time.Time
}

func NewOrder(orderID string) (*Order, error) {
	id, err := NewOrderID(orderID)
	if err != nil {
		return nil, err
	}
	return &Order{
		OrderID:    id,
		UploadedAt: time.Now(),
	}, nil
}

var ErrOrderID = errors.New("order code not valid")

var DefaultOrderID = OrderID{
	Value: "",
}

type OrderID struct {
	Value string
}

func (id OrderID) String() string {
	return id.Value
}

func NewOrderID(value string) (OrderID, error) {
	if err := validate(value); err != nil {
		return DefaultOrderID, err
	}
	return OrderID{value}, nil
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
