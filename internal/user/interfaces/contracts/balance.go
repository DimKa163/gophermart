package contracts

import "github.com/DimKa163/gophermart/internal/shared/types"

type BalanceResponse struct {
	Current   types.Decimal `json:"current"`
	Withdrawn types.Decimal `json:"withdrawn"`
}
