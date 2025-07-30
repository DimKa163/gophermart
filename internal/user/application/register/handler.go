package register

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
)

var ErrLoginAlreadyExists = errors.New("user already exists")

type RegisterHandler struct {
	unitOfWork  uow.UnitOfWork
	authService auth.AuthService
}
type RegisterCommand struct {
	Login    string
	Password string
}

func New(unitOfWork uow.UnitOfWork, authService auth.AuthService) *RegisterHandler {
	return &RegisterHandler{unitOfWork: unitOfWork, authService: authService}
}

func (h *RegisterHandler) Handle(ctx context.Context, command *RegisterCommand) (*types.AppResult[string], error) {
	pwd, salt, err := h.authService.GenerateHash([]byte(command.Password))
	if err != nil {
		return nil, err
	}
	user := model.NewUser(command.Login, pwd, salt)
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
	err = balanceRepository.Insert(ctx, &model.BonusBalance{UserID: id})
	if err != nil {
		_ = tuw.Rollback(ctx)
		return nil, err
	}

	user, _ = userRepository.Get(ctx, command.Login)
	token, err := h.authService.Authenticate(user.ID, []byte(command.Password), user.Password, user.Salt)
	if err != nil {
		_ = tuw.Rollback(ctx)
		return nil, err
	}
	if err = tuw.Commit(ctx); err != nil {
		return nil, err
	}
	return &types.AppResult[string]{Code: types.Created, Payload: token}, nil
}
