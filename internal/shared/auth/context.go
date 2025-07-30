package auth

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
)

var ErrUserNotFound = errors.New("user not found in context")

func User(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value("userID").(int64)
	if !ok {
		return -1, ErrUserNotFound
	}
	return userID, nil
}

func SetUser(ctx context.Context, userId int64) context.Context {
	gCtx, ok := ctx.(*gin.Context)
	if !ok {
		ctx = context.WithValue(ctx, "userID", userId)
		return ctx
	}
	gCtx.Set("userID", userId)
	return gCtx
}
