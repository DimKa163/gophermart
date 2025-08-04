package accrual

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/external/accrual/dto"
	"io"
	"net/http"
	"time"
)

var ErrNoContent = NoContentError{
	Message: "Order not found",
}

type AccrualClient interface {
	Order(ctx context.Context, number string) (*dto.Order, error)
}

type accrualClient struct {
	addr       string
	httpClient *http.Client
}

func (a accrualClient) Order(ctx context.Context, number string) (*dto.Order, error) {
	fullAddr := fmt.Sprintf("%s/api/order/%s", a.addr, number)
	req, err := http.NewRequestWithContext(ctx, "GET", fullAddr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNoContent {
			return nil, ErrNoContent
		}
		return nil, errors.New(resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var order dto.Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func New(addr string, transportFactories []func(transport http.RoundTripper) http.RoundTripper) AccrualClient {
	var transport http.RoundTripper
	defaultTransport := &http.Transport{}
	transport = defaultTransport
	for _, t := range transportFactories {
		transport = t(transport)
	}
	return &accrualClient{
		httpClient: &http.Client{
			Transport: transport,
			Timeout:   30 * time.Second,
		},
		addr: addr,
	}
}

type NoContentError struct {
	Message string
}

func (e NoContentError) Error() string {
	return e.Message
}
