package contracts

import (
	"github.com/DimKa163/gophermart/internal/shared/types"
	"time"
)

type OrderItem struct {
	Number     string        `json:"number"`
	Status     string        `json:"status"`
	Accrual    types.Decimal `json:"accrual,omitempty"`
	UploadedAt *time.Time    `json:"uploaded_at,omitempty"`
}
