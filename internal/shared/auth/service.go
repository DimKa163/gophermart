package auth

import (
	"crypto/rand"
	"errors"
	"golang.org/x/crypto/argon2"
)

var ErrInvalidPassword = errors.New("invalid password")

type ArgonConfig struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint32
	SaltLength  uint32
	KeyLength   uint32
}
type AuthService interface {
	GenerateHash(password []byte) (pwd, salt []byte, err error)

	Authenticate(userID int64, password, hashedPassword, salt []byte) (string, error)

	Verify(token string) (*Claims, error)
}

type argonAuthService struct {
	ArgonConfig
	engine *JWTEngine
}

func (a *argonAuthService) Verify(token string) (*Claims, error) {
	cl, err := a.engine.ReadToken(token)
	if err != nil {
		return nil, err
	}
	return cl, nil
}

func (a *argonAuthService) GenerateHash(password []byte) (pwd, salt []byte, err error) {
	salt, err = a.generateSalt()
	if err != nil {
		pwd = nil
		salt = nil
		return
	}
	pwd = a.hash(password, salt)
	return
}

func (a *argonAuthService) Authenticate(userID int64, password, hashedPassword, salt []byte) (string, error) {
	candidateHash := a.hash(password, salt)
	if !a.compare(hashedPassword, candidateHash) {
		return "", ErrInvalidPassword
	}
	return a.engine.BuildToken(userID)
}

func (a *argonAuthService) hash(pwd []byte, salt []byte) []byte {
	return argon2.IDKey([]byte(pwd), salt, a.Iterations, a.Memory, uint8(a.Parallelism), a.KeyLength)
}
func (a *argonAuthService) generateSalt() ([]byte, error) {
	salt := make([]byte, a.SaltLength)
	_, err := rand.Read(salt)
	return salt, err
}

func (a *argonAuthService) compare(b, c []byte) bool {
	if len(b) != len(c) {
		return false
	}
	var result byte
	for i := range b {
		result |= b[i] ^ c[i]
	}
	return result == 0
}

func NewAuthService(config ArgonConfig, jwt *JWTEngine) AuthService {
	return &argonAuthService{ArgonConfig: config, engine: jwt}
}
