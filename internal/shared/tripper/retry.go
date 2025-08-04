package tripper

import (
	"bytes"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/cenkalti/backoff/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type retryRoundTripper struct {
	rt http.RoundTripper
}

func NewRetryRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &retryRoundTripper{rt: rt}
}

func (rt *retryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte
	var err error
	if req.Body != nil {
		body, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewBuffer(body))
	}
	times := [3]int{1, 3, 5}
	attempt := 0
	logger := logging.Logger(req.Context())
	return backoff.Retry(req.Context(), func() (*http.Response, error) {

		response, err := rt.rt.RoundTrip(req)
		if err != nil {
			logger.Error("error occurred during request", zap.Error(err))
			return nil, backoff.Permanent(err)
		}

		if rt.shouldRetry(response) && attempt < len(times) {
			seconds := times[attempt]
			logger.Warn("try again...", zap.Int("attempt", attempt))
			if err = rt.drain(response); err != nil {

				logger.Error("error occurred during request", zap.Error(err))
				return nil, backoff.Permanent(err)
			}
			if req.Body != nil {
				req.Body = io.NopCloser(bytes.NewBuffer(body))
			}
			retryAfter := response.Header.Get("Retry-After")
			if retryAfter != "" {
				seconds, err = strconv.Atoi(retryAfter)
				if err != nil {
					seconds = times[attempt]
				}
			}
			attempt++
			return nil, backoff.RetryAfter(times[seconds])
		}
		return response, nil
	}, backoff.WithBackOff(backoff.NewExponentialBackOff()))
}

func (rt *retryRoundTripper) shouldRetry(resp *http.Response) bool {
	switch resp.StatusCode {
	case http.StatusRequestTimeout:
	case http.StatusTooManyRequests:
	case http.StatusBadGateway:
	case http.StatusGatewayTimeout:
	case http.StatusServiceUnavailable:
	case http.StatusInternalServerError:
		return true
	default:
		return false
	}
	return false
}

func (rt *retryRoundTripper) drain(response *http.Response) error {
	if response.Body != nil {
		_, err := io.Copy(io.Discard, response.Body)
		if err != nil {
			return err
		}
		defer response.Body.Close()
	}
	return nil
}
