package login

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/shared/types"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

type LoginHandler struct {
	unitOfWork uow.UnitOfWork
	jwtBuilder *auth.JWT
}
type LoginCommand struct {
	Login    string
	Password string
}

func New(unitOfWork uow.UnitOfWork, builder *auth.JWT) *LoginHandler {
	return &LoginHandler{unitOfWork: unitOfWork, jwtBuilder: builder}
}

func (h *LoginHandler) Handle(ctx context.Context, command *LoginCommand) (*types.AppResult[string], error) {
	rep := h.unitOfWork.UserRepository()
	user, err := rep.Get(ctx, command.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &types.AppResult[string]{
				Code:  types.Problem,
				Error: err,
			}, nil
		}
		return nil, err
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(command.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if bcrypt.CompareHashAndPassword([]byte(pwdHash), user.Password) == nil {
		return nil, ErrInvalidPassword
	}
	token, err := h.jwtBuilder.BuildJWT(user.ID)
	if err != nil {
		return nil, err
	}
	return &types.AppResult[string]{
		Code:    types.Created,
		Payload: token,
	}, nil
}
