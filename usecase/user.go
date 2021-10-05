package usecase

import "github.com/thak1411/rn-game-land-server/model"

type UserUsecase interface {
	CreateUser(model.User) error
}

type UserUS struct{}

func (us UserUS) CreateUser(user model.User) error {

	return nil
}

func NewUser() UserUsecase {
	return &UserUS{}
}
