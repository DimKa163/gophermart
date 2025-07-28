package types

import "github.com/shopspring/decimal"

type Decimal decimal.Decimal

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(decimal.Decimal(d).String()), nil
}

func NewDecimalFromString(str string) (Decimal, error) {
	dec, err := decimal.NewFromString(str)
	if err != nil {
		return Decimal{}, err
	}
	return Decimal(dec), nil
}
