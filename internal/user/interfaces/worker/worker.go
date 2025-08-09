package worker

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/user/application"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type OrderPooler struct {
	entryID  cron.EntryID
	signal   chan struct{}
	handler  *application.TrackOrderHandler
	schedule string
	limit    int
}

func NewWorker(cron *cron.Cron, schedule string, limit int, handler *application.TrackOrderHandler) (*OrderPooler, error) {
	signal := make(chan struct{})
	id, err := cron.AddFunc(schedule, func() {
		signal <- struct{}{}
	})
	if err != nil {
		return nil, err
	}
	return &OrderPooler{
		entryID:  id,
		signal:   signal,
		schedule: schedule,
		limit:    limit,
		handler:  handler,
	}, nil
}

func (w *OrderPooler) Run(ctx context.Context) error {

	logger := logging.Logger(ctx)
	ctx = logging.SetLogger(ctx, logger.With(zap.String("schedule", w.schedule)))
	go func() {
		for {
			select {
			case <-ctx.Done():
			case <-w.signal:
				logger := logging.Logger(ctx)
				logger.Info("start processing orders")
				err := w.handler.Handle(ctx, &application.TrackOrderCommand{
					Limit: 100,
				})
				if err != nil {
					logger.Warn("Failed to handle orders", zap.Error(err))
				}
			}
		}
	}()
	logger.Info("started")
	return nil
}
func (w *OrderPooler) Stop(ctx context.Context) {
	logger := logging.Logger(ctx)
	logger.Info("stopping processing orders")
	close(w.signal)
}
