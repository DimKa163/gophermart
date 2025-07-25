package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64
}

type JWTBuilderConfig struct {
	TokenExpiration time.Duration
	SecretKey       []byte
}

type JWTBuilder struct {
	JWTBuilderConfig
}

func NewJWTBuilder(config JWTBuilderConfig) *JWTBuilder {
	return &JWTBuilder{
		JWTBuilderConfig: config,
	}
}

func (b *JWTBuilder) BuildJWT(userID int64) (string, error) {
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
