package register

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"golang.org/x/crypto/bcrypt"
)

var ErrLoginAlreadyExists = errors.New("User already exists")

type RegisterHandler struct {
	unitOfWork uow.UnitOfWork
	jwtService *auth.JWT
}
type RegisterCommand struct {
	Login    string
	Password string
}

func New(unitOfWork uow.UnitOfWork, jwtService *auth.JWT) *RegisterHandler {
	return &RegisterHandler{unitOfWork: unitOfWork, jwtService: jwtService}
}

func (h *RegisterHandler) Handle(ctx context.Context, command *RegisterCommand) (*types.AppResult[string], error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := model.NewUser(command.Login, pwd)
	tuw, err := h.unitOfWork.Begin(ctx)
	if err != nil {
		return nil, err
	}
	userRepository := tuw.UserRepository()
	loginExists, err := userRepository.LoginExists(ctx, command.Login)
	if err != nil {
		return nil, err
	}
	if loginExists {
		return &types.AppResult[string]{
			Code:  types.Duplicate,
			Error: ErrLoginAlreadyExists,
		}, nil
	}

	id, err := userRepository.Insert(ctx, user)
	if err != nil {
		_ = tuw.Rollback(ctx)
		return nil, err
	}

	balanceRepository := tuw.BonusBalanceRepository()
	err = balanceRepository.Insert(ctx, &model.BonusBalance{UserId: id})
	if err != nil {
		_ = tuw.Rollback(ctx)
		return nil, err
	}

	token, err := h.jwtService.BuildJWT(id)
	if err != nil {
		_ = tuw.Rollback(ctx)
		return nil, err
	}
	if err = tuw.Commit(ctx); err != nil {
		return nil, err
	}
	return &types.AppResult[string]{Code: types.Created, Payload: token}, nil
}
