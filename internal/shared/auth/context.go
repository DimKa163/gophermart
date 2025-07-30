package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrUserNotFound = errors.New("user not found in context")

type UserID string

const (
	user UserID = "userID"
)

func User(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(user).(int64)
	if !ok {
		return -1, ErrUserNotFound
	}
	return userID, nil
}

func SetUser(ctx context.Context, userID int64) context.Context {
	gCtx, ok := ctx.(*gin.Context)
	if !ok {
		ctx = context.WithValue(ctx, user, userID)
		return ctx
	}
	gCtx.Set(string(user), userID)
	return gCtx
}
