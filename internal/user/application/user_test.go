package application

import (
	"context"
	"github.com/DimKa163/gophermart/internal/shared/auth"
	"github.com/DimKa163/gophermart/internal/user/domain/model"
	"github.com/DimKa163/gophermart/internal/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterShouldSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	login := "login"
	password := "password"
	salt := []byte("salt")
	passwordHash := []byte(password)
	token := "token"
	us := &model.User{
		Login:    login,
		Password: passwordHash,
		Salt:     salt,
	}
	mockUow.EXPECT().UserRepository().Return(mockRepo)

	mockRepo.EXPECT().LoginExists(ctx, login).Return(false, nil)

	mockAuth.EXPECT().GenerateHash([]byte(password)).Return(passwordHash, salt, nil)

	mockRepo.EXPECT().Insert(ctx, us).Return(int64(1), nil)

	mockRepo.EXPECT().Get(ctx, login).Return(us, nil)

	mockAuth.EXPECT().Authenticate(us.ID, []byte(password), us.Password, us.Salt).Return(token, nil)

	sut := NewUserService(mockUow, mockAuth)

	result, err := sut.Register(ctx, login, password)

	assert.NoError(t, err, "Register should succeed")
	assert.Equal(t, token, result, "token should match")
}

func TestRegisterWithKnownLoginShouldFailer(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	login := "login"
	password := "password"

	token := ""

	mockUow.EXPECT().UserRepository().Return(mockRepo)

	mockRepo.EXPECT().LoginExists(ctx, login).Return(true, nil)

	sut := NewUserService(mockUow, mockAuth)

	result, err := sut.Register(ctx, login, password)

	assert.ErrorIs(t, ErrLoginAlreadyExists, err, "Register should fail")
	assert.Equal(t, token, result, "token should match")
}

func TestLoginShouldSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	login := "login"
	password := "password"
	salt := []byte("salt")
	passwordHash := []byte(password)
	token := "token"
	us := &model.User{
		ID:       1,
		Login:    login,
		Password: passwordHash,
		Salt:     salt,
	}

	mockUow.EXPECT().UserRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, login).Return(us, nil)

	mockAuth.EXPECT().Authenticate(us.ID, []byte(password), us.Password, us.Salt).Return(token, nil)

	sut := NewUserService(mockUow, mockAuth)

	result, err := sut.Login(ctx, login, password)

	assert.NoError(t, err, "Login should succeed")
	assert.Equal(t, token, result, "token should match")
}

func TestLoginWithWrongPwdShouldFail(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	login := "login"
	password := "password"
	salt := []byte("salt")
	passwordHash := []byte(password)
	token := ""
	us := &model.User{
		ID:       1,
		Login:    login,
		Password: passwordHash,
		Salt:     salt,
	}
	mockUow.EXPECT().UserRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, login).Return(us, nil)

	mockAuth.EXPECT().Authenticate(us.ID, []byte(password), us.Password, us.Salt).Return("", auth.ErrInvalidPassword)

	sut := NewUserService(mockUow, mockAuth)

	result, err := sut.Login(ctx, login, password)

	assert.ErrorIs(t, auth.ErrInvalidPassword, err, "Login should fail")
	assert.Equal(t, token, result, "token should match")
}

func TestLoginWithWrongLoginShouldFail(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUow := mocks.NewMockUnitOfWork(ctrl)
	mockAuth := mocks.NewMockAuthService(ctrl)
	mockRepo := mocks.NewMockUserRepository(ctrl)

	login := "login"
	password := "password"
	token := ""
	mockUow.EXPECT().UserRepository().Return(mockRepo)

	mockRepo.EXPECT().Get(ctx, login).Return(nil, pgx.ErrNoRows)
	sut := NewUserService(mockUow, mockAuth)
	result, err := sut.Login(ctx, login, password)

	assert.ErrorIs(t, ErrUserNotFound, err, "Login should fail")
	assert.Equal(t, token, result, "token should match")
}
