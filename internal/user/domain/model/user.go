package model

import (
	"github.com/jackc/pgx/v5"
	"time"
)

type User struct {
	ID        int64
	CreatedAt time.Time
	Login     string
	Password  []byte
}

func NewUser(login string, password []byte) *User {
	return &User{
		CreatedAt: time.Now(),
		Login:     login,
		Password:  password,
	}
}

func (u *User) Scan(row pgx.Row) error {
	if err := row.Scan(&u.ID, &u.CreatedAt, &u.Login, &u.Password); err != nil {
		return err
	}
	return nil
}
