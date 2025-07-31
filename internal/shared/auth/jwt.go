package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

type JWTConfig struct {
	TokenExpiration time.Duration
	SecretKey       []byte
}

type JWTEngine struct {
	JWTConfig
}

func NewJWT(config JWTConfig) *JWTEngine {
	return &JWTEngine{
		JWTConfig: config,
	}
}

func (b *JWTConfig) BuildToken(userID int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(b.TokenExpiration)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString(b.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (b *JWTConfig) ReadToken(tokenString string) (*Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(strings.ReplaceAll(tokenString, "Bearer ", ""), &claims, func(token *jwt.Token) (interface{}, error) {
		return b.SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return &claims, nil
}
