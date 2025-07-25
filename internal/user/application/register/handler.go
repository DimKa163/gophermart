package register

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"golang.org/x/crypto/bcrypt"
)

var ErrLoginAlreadyExists = errors.New("User already exists")

type RegisterHandler struct {
	unitOfWork uow.UnitOfWork
}

func New(unitOfWork uow.UnitOfWork) *RegisterHandler {
	return &RegisterHandler{unitOfWork: unitOfWork}
}

func (h *RegisterHandler) Handle(ctx context.Context, command *RegisterCommand) (any, error) {
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
		return nil, ErrLoginAlreadyExists
	}
	id, err := userRepository.Insert(ctx, user)
	if err != nil {
		_ = tuw.Rollback(ctx)
		return nil, err
	}

	if err = tuw.Commit(ctx); err != nil {
		return nil, err
	}
	return id, nil
}
