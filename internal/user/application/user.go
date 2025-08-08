package application

import (
	"context"
	"errors"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/user/domain"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/domain/uow"
	"github.com/jackc/pgx/v5"
)

var (
	ErrLoginAlreadyExists = domain.NewLoginAlreadyExists("user already exists")
	ErrUserNotFound       = domain.NewResourceNotFound("user not found")
)

type userService struct {
	uow  uow.UnitOfWork
	auth auth.AuthService
}

func (u *userService) Register(ctx context.Context, login string, password string) (string, error) {
	pwd, salt, err := u.auth.GenerateHash([]byte(password))
	if err != nil {
		return "", err
	}
	user := model.NewUser(login, pwd, salt)
	userRep := u.uow.UserRepository()
	loginExists, err := userRep.LoginExists(ctx, login)
	if err != nil {
		return "", err
	}
	if loginExists {
		return "", ErrLoginAlreadyExists
	}
	_, err = userRep.Insert(ctx, user)
	if err != nil {
		return "", err
	}

	user, _ = userRep.Get(ctx, login)
	token, err := u.authenticate(user, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) Login(ctx context.Context, login string, password string) (string, error) {
	userRep := u.uow.UserRepository()
	user, err := userRep.Get(ctx, login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUserNotFound
		}
		return "", err
	}
	token, err := u.authenticate(user, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *userService) Balance(ctx context.Context) (*model.BonusBalance, error) {
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	bal, err := u.uow.BonusBalanceRepository().Get(ctx, userID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if bal == nil {
		bal = &model.BonusBalance{}
	}
	return bal, nil
}

func (u *userService) Withdrawal(ctx context.Context) ([]*model.Transaction, error) {
	userID, err := auth.User(ctx)
	if err != nil {
		return nil, err
	}
	rep := u.uow.BonusMovementRepository()
	t := model.WITHDRAWAL
	items, err := rep.GetAll(ctx, userID, &t)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (u *userService) authenticate(user *model.User, password string) (string, error) {
	token, err := u.auth.Authenticate(user.ID, []byte(password), user.Password, user.Salt)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidPassword) {
			return "", err
		}
		return "", err
	}
	return token, nil
}

func NewUserService(uow uow.UnitOfWork, auth auth.AuthService) domain.UserService {
	return &userService{
		uow:  uow,
		auth: auth,
	}
}
