package rest

import (
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/application/balance"
	"github.com/DimKa163/gophermart/internal/user/application/login"
	"github.com/DimKa163/gophermart/internal/user/application/order"
	"github.com/DimKa163/gophermart/internal/user/application/register"
	"github.com/DimKa163/gophermart/internal/user/application/withdrawal"
	"github.com/DimKa163/gophermart/internal/user/interfaces/contracts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type UserAPI interface {
	Register(context *gin.Context)
	Login(context *gin.Context)
	Upload(context *gin.Context)
	GetOrders(context *gin.Context)
	GetBalance(context *gin.Context)
	Withdraw(context *gin.Context)
	GetWithdrawals(context *gin.Context)
}

type userApi struct {
	registerHandler        *register.RegisterHandler
	loginHandler           *login.LoginHandler
	uploadHandler          *order.UploadOrderHandler
	orderQueryHandler      *order.OrderQueryHandler
	bonusBalanceHandler    *balance.BalanceQueryHandler
	withdrawHandler        *balance.WithdrawHandler
	withdrawalQueryHandler *withdrawal.WithdrawalQueryHandler
}

func NewUserApi(
	registerHandler *register.RegisterHandler,
	loginHandler *login.LoginHandler,
	uploadOrderHandler *order.UploadOrderHandler,
	orderQueryHandler *order.OrderQueryHandler,
	bonusBalanceHandler *balance.BalanceQueryHandler,
	withdrawHandler *balance.WithdrawHandler,
	withdrawalQueryHandler *withdrawal.WithdrawalQueryHandler,
) UserAPI {
	return &userApi{
		registerHandler:        registerHandler,
		loginHandler:           loginHandler,
		uploadHandler:          uploadOrderHandler,
		orderQueryHandler:      orderQueryHandler,
		bonusBalanceHandler:    bonusBalanceHandler,
		withdrawHandler:        withdrawHandler,
		withdrawalQueryHandler: withdrawalQueryHandler,
	}
}

func (u *userApi) Register(context *gin.Context) {
	logger := logging.Logger(context)
	var user contracts.User
	if err := context.ShouldBind(&user); err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.registerHandler.Handle(context, &register.RegisterCommand{Login: user.Login, Password: user.Password})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.Created:
		context.Header("Authorization", result.Payload)
		context.Status(http.StatusOK)
		break
	case types.Duplicate:
		logger.Error("Error occurred", zap.Error(result.Error))
		context.JSON(http.StatusConflict, gin.H{"error": result.Error})
		break
	}
}

func (u *userApi) Login(context *gin.Context) {
	logger := logging.Logger(context)
	var user contracts.User
	if err := context.ShouldBind(&user); err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.loginHandler.Handle(context, &login.LoginCommand{Login: user.Login, Password: user.Password})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.Created:
		context.Header("Authorization", result.Payload)
		context.Status(http.StatusOK)
		break
	case types.Problem:
		logger.Error("Error occurred", zap.Error(err))
		context.JSON(http.StatusUnauthorized, gin.H{"error": result.Error.Error()})
		break
	}
}

func (u *userApi) Upload(context *gin.Context) {
	logger := logging.Logger(context)
	body, err := context.GetRawData()
	if err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.uploadHandler.Handle(context, &order.UploadOrderCommand{Id: string(body)})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.Created:
		context.Status(http.StatusAccepted)
		break
	case types.NoChange:
		context.Status(http.StatusOK)
		break
	case types.Problem:
		logger.Error("Error occurred", zap.Error(result.Error))
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": result.Error})
	case types.Duplicate:
		logger.Error("Error occurred", zap.Error(result.Error))
		context.JSON(http.StatusConflict, gin.H{"error": result.Error})
		break
	}
}

func (u *userApi) GetOrders(context *gin.Context) {
	logger := logging.Logger(context)
	result, err := u.orderQueryHandler.Handle(context, &order.OrderQuery{})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.NoChange:
		if len(result.Payload) == 0 {
			context.Status(http.StatusNoContent)
			return
		}
		orderItems := make([]contracts.OrderItem, len(result.Payload))
		for i, item := range result.Payload {
			orderItems[i] = contracts.OrderItem{
				Number:     item.OrderID.String(),
				Accrual:    item.Accrual,
				Status:     item.Status.String(),
				UploadedAt: item.UploadedAt,
			}
		}
		context.JSON(http.StatusOK, orderItems)
		break
	}
}

func (u *userApi) GetBalance(context *gin.Context) {
	logger := logging.Logger(context)
	result, err := u.bonusBalanceHandler.Handle(context, &balance.BalanceQuery{})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.NoChange:
		context.JSON(http.StatusOK, contracts.BalanceResponse{Current: result.Payload.Current, Withdrawn: result.Payload.Withdrawn})
	}
}

func (u *userApi) Withdraw(context *gin.Context) {
	logger := logging.Logger(context)
	var body contracts.WithdrawRequest
	if err := context.ShouldBind(&body); err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.withdrawHandler.Handle(context, &balance.WithdrawCommand{
		OrderID: body.OrderID,
		Sum:     body.Sum,
	})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.Created:
		context.Status(http.StatusOK)
		return
	case types.Problem:
		logger.Error("Error occurred", zap.Error(result.Error))
		if errors.Is(result.Error, balance.ErrNegativeBalance) {
			context.JSON(http.StatusPaymentRequired, gin.H{"error": result.Error.Error()})
			return
		}
		if errors.Is(result.Error, balance.ErrWrongOrder) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{"error": result.Error.Error()})
			return
		}
	}
}

func (u *userApi) GetWithdrawals(context *gin.Context) {
	logger := logging.Logger(context)
	result, err := u.withdrawalQueryHandler.Handle(context, &withdrawal.WithdrawalQuery{})
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	switch result.Code {
	case types.NoChange:
		if len(result.Payload) == 0 {
			context.Status(http.StatusNoContent)
			return
		}
		response := make([]contracts.WithdrawResponse, len(result.Payload))
		for i, item := range result.Payload {
			response[i] = contracts.WithdrawResponse{
				OrderID:     item.OrderID,
				Sum:         item.Amount,
				ProcessedAt: item.CreatedAt,
			}
		}
		context.JSON(http.StatusOK, response)
		break
	}
}
