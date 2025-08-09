package application

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUploadNewOrderShouldSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockRepo := mocks.NewMockOrderRepository(ctrl)

	ctx = auth.SetUser(ctx, 1)
	number := "12345678903"
	orderID, _ := model.NewOrderID(number)
	ord := &model.Order{
		OrderID: orderID,
		UserID:  1,
		Status:  model.OrderStatusNEW,
	}
	mockUow.EXPECT().OrderRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, orderID).Return(nil, pgx.ErrNoRows)

	mockRepo.EXPECT().Insert(ctx, ord).Return(orderID, nil)

	sut := NewOrderService(mockUow)

	result, err := sut.Upload(ctx, orderID)

	assert.NoError(t, err, "Upload should return no error")
	assert.True(t, result, "Upload should return result true")
}

func TestUploadNotNewOrderShouldReturnSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockRepo := mocks.NewMockOrderRepository(ctrl)

	ctx = auth.SetUser(ctx, 1)
	number := "12345678903"
	orderID, _ := model.NewOrderID(number)
	ord := &model.Order{
		OrderID: orderID,
		UserID:  1,
		Status:  model.OrderStatusNEW,
	}
	mockUow.EXPECT().OrderRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, orderID).Return(ord, nil)

	sut := NewOrderService(mockUow)

	result, err := sut.Upload(ctx, orderID)

	assert.NoError(t, err, "Upload should return no error")
	assert.False(t, result, "Upload should return result true")
}

func TestUploadNotNewOrderShouldReturnError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockRepo := mocks.NewMockOrderRepository(ctrl)

	ctx = auth.SetUser(ctx, 1)
	number := "12345678903"
	orderID, _ := model.NewOrderID(number)
	ord := &model.Order{
		OrderID: orderID,
		UserID:  2,
		Status:  model.OrderStatusNEW,
	}
	mockUow.EXPECT().OrderRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, orderID).Return(ord, nil)

	sut := NewOrderService(mockUow)

	result, err := sut.Upload(ctx, orderID)

	assert.ErrorIs(t, ErrOrderExistsWithAnotherUser, err, "Upload should return error for another user")
	assert.False(t, result, "Upload should return result true")
}

func TestWithdrawWithBalanceShouldSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockRepo := mocks.NewMockOrderRepository(ctrl)
	mockURepo := mocks.NewMockUserRepository(ctrl)
	ctx = auth.SetUser(ctx, 1)
	number := "12345678903"
	orderID, _ := model.NewOrderID(number)
	ord := &model.Order{
		OrderID: orderID,
		UserID:  1,
		Status:  model.OrderStatusNEW,
	}
	ord.AddTransaction(model.WITHDRAWAL, types.Decimal{Decimal: decimal.NewFromFloat32(100.00)})

	bal := &model.BonusBalance{UserID: 1, Current: types.Decimal{Decimal: decimal.NewFromFloat32(500.00)}, Withdrawn: types.Decimal{Decimal: decimal.NewFromFloat32(0.00)}}

	mockUow.EXPECT().OrderRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, orderID).Return(ord, nil)

	mockUow.EXPECT().UserRepository().Return(mockURepo)

	mockURepo.EXPECT().GetBonusBalanceByUserID(ctx, int64(1)).Return(bal, nil)

	mockRepo.EXPECT().Update(ctx, ord).Return(nil)

	sut := NewOrderService(mockUow)

	err := sut.Withdraw(ctx, orderID, types.Decimal{Decimal: decimal.NewFromFloat32(100.00)})

	assert.NoError(t, err, "Withdraw should return no error")
}

func TestWithdrawWithoutBalanceShouldReturnError(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockRepo := mocks.NewMockOrderRepository(ctrl)
	mockURepo := mocks.NewMockUserRepository(ctrl)
	ctx = auth.SetUser(ctx, 1)
	number := "12345678903"
	orderID, _ := model.NewOrderID(number)
	ord := &model.Order{
		OrderID: orderID,
		UserID:  1,
		Status:  model.OrderStatusNEW,
	}
	ord.AddTransaction(model.WITHDRAWAL, types.Decimal{Decimal: decimal.NewFromFloat32(100.00)})

	bal := &model.BonusBalance{UserID: 1, Current: types.Decimal{Decimal: decimal.NewFromFloat32(50.00)}, Withdrawn: types.Decimal{Decimal: decimal.NewFromFloat32(0.00)}}

	mockUow.EXPECT().OrderRepository().Return(mockRepo).MaxTimes(2)

	mockRepo.EXPECT().Get(ctx, orderID).Return(ord, nil).MaxTimes(2)

	mockUow.EXPECT().UserRepository().Return(mockURepo)

	mockURepo.EXPECT().GetBonusBalanceByUserID(ctx, int64(1)).Return(bal, nil)

	sut := NewOrderService(mockUow)

	err := sut.Withdraw(ctx, orderID, types.Decimal{Decimal: decimal.NewFromFloat32(100.00)})

	assert.ErrorIs(t, ErrNegativeBalance, err, "Withdraw should return error")
}
