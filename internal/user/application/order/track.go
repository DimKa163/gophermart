package order

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/external/accrual"
	"github.com/DimKa163/gophermart/internal/user/infrastructure/external/accrual/dto"
	"sync"
)

var statusMap map[dto.OrderStatus]model.OrderStatus = map[dto.OrderStatus]model.OrderStatus{
	dto.StatusREGISTERED: model.OrderStatusNEW,
	dto.StatusPROCESSING: model.OrderStatusPROCESSING,
	dto.StatusPROCESSED:  model.OrderStatusPROCESSED,
	dto.StatusINVALID:    model.OrderStatusINVALID,
}

type TrackOrderCommand struct {
	Limit int
}

type TrackOrderHandler struct {
	uow       uow.UnitOfWork
	processor *TrackOrderProcessor
}

func NewTrackOrderHandler(uow uow.UnitOfWork, processor *TrackOrderProcessor) *TrackOrderHandler {
	return &TrackOrderHandler{uow: uow, processor: processor}
}

func (handler *TrackOrderHandler) Handle(ctx context.Context, command *TrackOrderCommand) (*types.AppResult[any], error) {
	var err error
	txUow, err := handler.uow.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = txUow.Rollback(ctx)
			return
		}
	}()
	orderRep := txUow.OrderRepository()
	offset := 0
	items, err := orderRep.GetForUpdate(ctx, command.Limit, offset, model.OrderStatusNEW, model.OrderStatusPROCESSING)
	if err != nil {
		return nil, err
	}
	for len(items) > 0 {
		ch := handler.processor.Process(ctx, items)
		for it := range ch {
			err = orderRep.Update(ctx, it)
		}
		offset += command.Limit
		items, err = orderRep.GetForUpdate(ctx, command.Limit, offset, model.OrderStatusNEW, model.OrderStatusPROCESSING)
		if err != nil {
			return nil, err
		}
	}
	_ = txUow.Commit(ctx)
	return &types.AppResult[any]{}, nil
}

type TrackOrderProcessor struct {
	accrualCl accrual.AccrualClient
}

func (p *TrackOrderProcessor) Process(ctx context.Context, orders []*model.Order) <-chan *model.Order {
	inputCh := p.iterate(ctx, orders)
	channels := p.fanOut(ctx, len(orders), inputCh, p.processOrder)
	return p.fanIn(ctx, channels...)
}

func (p *TrackOrderProcessor) iterate(ctx context.Context, orders []*model.Order) <-chan *model.Order {
	inputCh := make(chan *model.Order)
	go func() {
		defer close(inputCh)
		for _, order := range orders {
			select {
			case <-ctx.Done():
				return
			case inputCh <- order:
			}
		}
	}()
	return inputCh
}

func (p *TrackOrderProcessor) fanOut(ctx context.Context, len int, inputCh <-chan *model.Order,
	fn func(ctx context.Context,
		inputCh <-chan *model.Order) <-chan *model.Order) []<-chan *model.Order {
	channels := make([]<-chan *model.Order, len)
	for i := 0; i < len; i++ {
		channels[i] = fn(ctx, inputCh)
	}
	return channels
}

func (p *TrackOrderProcessor) fanIn(ctx context.Context, channels ...<-chan *model.Order) <-chan *model.Order {
	finalCh := make(chan *model.Order)
	var wg sync.WaitGroup
	for _, ch := range channels {
		chClosure := ch
		wg.Add(1)
		go func() {
			defer wg.Done()
			for data := range chClosure {
				select {
				case <-ctx.Done():
					return
				case finalCh <- data:
				}
			}
		}()

	}
	go func() {
		wg.Wait()
		close(finalCh)
	}()
	return finalCh
}

func (p *TrackOrderProcessor) processOrder(ctx context.Context, inputCh <-chan *model.Order) <-chan *model.Order {
	info := make(chan *model.Order)
	go func() {
		defer close(info)
		for data := range inputCh {
			if data == nil {
				continue
			}
			data.Error = ""
			or, err := p.accrualCl.Order(ctx, data.OrderID.String())
			if err != nil {
				data.Error = err.Error()
			} else {
				if data.Status != statusMap[or.Status] {
					data.Status = statusMap[or.Status]
					if or.Accrual != nil && !or.Accrual.IsZero() && data.Status == model.OrderStatusPROCESSED {
						data.AddTransaction(model.ACCRUAL, *or.Accrual)
					}
				}
			}

			select {
			case <-ctx.Done():
				return
			case info <- data:
			}
		}
	}()
	return info
}

func NewTrackOrderProcessor(accrualCl accrual.AccrualClient) *TrackOrderProcessor {
	return &TrackOrderProcessor{accrualCl: accrualCl}
}
