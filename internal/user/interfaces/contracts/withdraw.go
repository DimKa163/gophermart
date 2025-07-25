package contracts

import "time"

type WithdrawRequest struct {
	OrderID string  `json:"order"`
	Sum     float64 `json:"sum"`
}

type WithdrawResponse struct {
	OrderID     string    `json:"order"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
