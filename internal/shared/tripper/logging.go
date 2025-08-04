package tripper

import (
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"go.uber.org/zap"
	"net/http"
)

type loggingRoundTripper struct {
	rt http.RoundTripper
}

func (s *loggingRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	logger := logging.Logger(r.Context()).With(zap.String("method", r.Method),
		zap.String("url", r.URL.String()))
	r = r.WithContext(logging.SetLogger(r.Context(), logger))
	logger.Info("making request")
	resp, err := s.rt.RoundTrip(r)
	if err != nil {
		logger.Error("error occurred during request", zap.Error(err))
		return nil, err
	}
	logger = logging.Logger(r.Context()).With(zap.String("status", resp.Status))
	logger.Info("response received")
	return resp, nil
}

func NewLoggingRoundTripper(rt http.RoundTripper) http.RoundTripper {
	return &loggingRoundTripper{rt}
}
