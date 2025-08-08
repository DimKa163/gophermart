package rest

import (
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/logging"
	"github.com/DimKa163/gophermart/internal/user/application"
	"github.com/DimKa163/gophermart/internal/user/domain"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
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

type userAPI struct {
	user  domain.UserService
	order domain.OrderService
}

func NewUserAPI(user domain.UserService, order domain.OrderService) UserAPI {
	return &userAPI{
		user:  user,
		order: order,
	}
}

func (u *userAPI) Register(context *gin.Context) {
	logger := logging.Logger(context)
	var user contracts.User
	if err := context.ShouldBind(&user); err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.user.Register(context, user.Login, user.Password)
	if err != nil {
		if errors.Is(application.ErrLoginAlreadyExists, err) {
			context.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		logger.Error("unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Header("Authorization", result)
	context.Status(http.StatusOK)
}

func (u *userAPI) Login(context *gin.Context) {
	logger := logging.Logger(context)
	var user contracts.User
	if err := context.ShouldBind(&user); err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := u.user.Login(context, user.Login, user.Password)
	if err != nil {
		if errors.Is(application.ErrUserNotFound, err) {
			context.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		logger.Error("unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Header("Authorization", result)
	context.Status(http.StatusOK)
}

func (u *userAPI) Upload(context *gin.Context) {
	logger := logging.Logger(context)
	body, err := context.GetRawData()
	if err != nil {
		logger.Error("Error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := u.order.Upload(context, string(body))
	if err != nil {
		if errors.Is(model.ErrOrderID, err) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(application.ErrOrderExistsWithAnotherUser, err) {
			context.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		logger.Error("unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !result {
		context.Status(http.StatusOK)
		return
	}
	context.Status(http.StatusAccepted)
}

func (u *userAPI) GetOrders(context *gin.Context) {
	logger := logging.Logger(context)
	result, err := u.order.List(context)
	if err != nil {
		logger.Error("unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(result) == 0 {
		context.Status(http.StatusNoContent)
		return
	}
	orderItems := make([]contracts.OrderItem, len(result))
	for i, item := range result {
		orderItems[i] = contracts.OrderItem{
			Number:     item.OrderID.String(),
			Accrual:    item.Accrual,
			Status:     item.Status.String(),
			UploadedAt: item.UploadedAt,
		}
	}
	context.JSON(http.StatusOK, orderItems)
}

func (u *userAPI) GetBalance(context *gin.Context) {
	logger := logging.Logger(context)
	result, err := u.user.Balance(context)
	if err != nil {
		logger.Error("unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, contracts.BalanceResponse{Current: &result.Current, Withdrawn: &result.Withdrawn})
}

func (u *userAPI) Withdraw(context *gin.Context) {
	logger := logging.Logger(context)
	var body contracts.WithdrawRequest
	if err := context.ShouldBind(&body); err != nil {
		logger.Error("error reading body", zap.Error(err))
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := u.order.Withdraw(context, body.OrderID, body.Sum)
	if err != nil {
		if errors.Is(model.ErrOrderID, err) {
			context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, application.ErrNegativeBalance) {
			context.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
			return
		}
		logger.Error("unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.Status(http.StatusOK)
}

func (u *userAPI) GetWithdrawals(context *gin.Context) {
	logger := logging.Logger(context)
	result, err := u.user.Withdrawal(context)
	if err != nil {
		logger.Error("Unhandled error occurred", zap.Error(err))
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(result) == 0 {
		context.Status(http.StatusNoContent)
		return
	}
	response := make([]contracts.WithdrawResponse, len(result))
	for i, item := range result {
		response[i] = contracts.WithdrawResponse{
			OrderID:     item.OrderID,
			Sum:         item.Amount,
			ProcessedAt: item.CreatedAt,
		}
	}
	context.JSON(http.StatusOK, response)
}
