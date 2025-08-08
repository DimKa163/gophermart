package contracts

import (
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"time"
)

type WithdrawRequest struct {
	OrderID string        `json:"order"`
	Sum     types.Decimal `json:"sum"`
}

type WithdrawResponse struct {
	OrderID     model.OrderID `json:"order"`
	Sum         types.Decimal `json:"sum"`
	ProcessedAt time.Time     `json:"processed_at"`
}
