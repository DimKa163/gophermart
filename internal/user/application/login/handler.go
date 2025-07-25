package login

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrInvalidPassword = errors.New("invalid password")

type LoginHandler struct {
	unitOfWork uow.UnitOfWork
	jwtBuilder *auth.JWTBuilder
}

func New(unitOfWork uow.UnitOfWork, builder *auth.JWTBuilder) *LoginHandler {
	return &LoginHandler{unitOfWork: unitOfWork, jwtBuilder: builder}
}

func (h *LoginHandler) Handle(ctx context.Context, command *LoginCommand) (any, error) {
	rep := h.unitOfWork.UserRepository()
	user, err := rep.Get(ctx, command.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
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
	return token, nil
}
