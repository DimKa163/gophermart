package types

import (
	"database/sql/driver"
	"fmt"
	"github.com/shopspring/decimal"
)

type Decimal struct {
	decimal.Decimal
}

func (d *Decimal) MarshalJSON() ([]byte, error) {

	a := []byte(d.String())
	return a, nil
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	strVal := string(data)
	dec, err := decimal.NewFromString(strVal)
	if err != nil {
		return fmt.Errorf("Decimal.UnmarshalJSON: invalid input: %s", string(data))
	}
	d.Decimal = dec
	return nil
}

func NewDecimalFromString(str string) (Decimal, error) {
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return Decimal{}, err
	}
	return Decimal{dec}, nil
}

func (d *Decimal) Cmp(a Decimal) int {
	return d.Decimal.Cmp(a.Decimal)
}

func (d *Decimal) Add(a Decimal) Decimal {
	return Decimal{d.Decimal.Add(a.Decimal)}
}

func (d *Decimal) Sub(a Decimal) Decimal {
	return Decimal{d.Decimal.Sub(a.Decimal)}
}

func (d *Decimal) IsNegative() bool {
	return d.Decimal.IsNegative()
}

func (d *Decimal) IsPositive() bool {
	return d.Decimal.IsPositive()
}
func (d *Decimal) IsZero() bool {
	return d.Decimal.IsZero()
}
func (d *Decimal) Value() (driver.Value, error) {
	return d.Decimal.String(), nil
}
