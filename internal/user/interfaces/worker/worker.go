package worker

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/shared/mediatr"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/application/order"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Worker struct {
	cron     *cron.Cron
	schedule string
}

func NewWorker(cron *cron.Cron, schedule string) *Worker {
	return &Worker{
		cron:     cron,
		schedule: schedule,
	}
}

func (w *Worker) Start(ctx context.Context) error {
	_, err := w.cron.AddFunc(w.schedule, func() {
		if ctx.Err() != nil {
			return
		}
		logger := logging.Logger(ctx)
		err := w.Run(ctx)
		if err != nil {
			logger.Warn("Failed to run job", zap.Error(err))
		}
	})
	if err != nil {
		return err
	}

	w.cron.Start()
	logger := logging.Logger(ctx)
	logger.Info("started")
	return nil
}

func (w *Worker) Run(ctx context.Context) error {
	logger := logging.Logger(ctx)
	logger.Info("start processing orders")
	_, err := mediatr.Send[*order.TrackOrderCommand, *types.AppResult[any]](ctx, &order.TrackOrderCommand{
		Limit: 100,
	})
	if err != nil {
		logger.Warn("Failed to handle orders", zap.Error(err))
		return err
	}
	return nil
}
