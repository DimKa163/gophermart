package model

import "errors"

var ErrBonusBalance = errors.New("invalid bonus balance")

type BonusBalance struct {
	UserId    int64
	Current   float64
	Withdrawn float64
}

func NewBonusBalance(userId int64, current, withdrawn float64) (*BonusBalance, error) {
	if current < 0 || withdrawn < 0 {
		return nil, ErrBonusBalance
	}
	return &BonusBalance{
		UserId:    userId,
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}
