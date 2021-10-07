package usecase

import (
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/model"
)

type UserUsecase interface {
	CreateUser(model.User) error
	UpdateUser(model.User) error
	DeleteUser(int) error
	GetAllUser() ([]model.User, error)
}

type UserUC struct {
	db database.UserDatabase
}

func (uc *UserUC) CreateUser(user model.User) error {
	// Inject Salt & Hasing Password //
	return uc.db.Create(user)
}

func (uc *UserUC) UpdateUser(user model.User) error {
	// Hashing Password //
	return uc.db.Update(user)
}

func (uc *UserUC) DeleteUser(id int) error {
	return uc.db.Delete(id)
}

func (uc *UserUC) GetAllUser() ([]model.User, error) {
	return uc.db.GetAll()
}

func NewUser(db database.UserDatabase) UserUsecase {
	return &UserUC{db}
}
