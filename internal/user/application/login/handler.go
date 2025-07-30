package login

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

var ErrUserNotFound = errors.New("user not found")

type LoginHandler struct {
	unitOfWork  uow.UnitOfWork
	authService auth.AuthService
}
type LoginCommand struct {
	Login    string
	Password string
}

func New(unitOfWork uow.UnitOfWork, authService auth.AuthService) *LoginHandler {
	return &LoginHandler{unitOfWork: unitOfWork, authService: authService}
}

func (h *LoginHandler) Handle(ctx context.Context, command *LoginCommand) (*types.AppResult[string], error) {
	rep := h.unitOfWork.UserRepository()
	user, err := rep.Get(ctx, command.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &types.AppResult[string]{
				Code:  types.Problem,
				Error: ErrUserNotFound,
			}, nil
		}
		return nil, err
	}
	token, err := h.authService.Authenticate(user.ID, []byte(command.Password), user.Password, user.Salt)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidPassword) {
			return &types.AppResult[string]{
				Code:  types.Problem,
				Error: err,
			}, nil
		}
		return nil, err
	}
	return &types.AppResult[string]{
		Code:    types.Created,
		Payload: token,
	}, nil
}
