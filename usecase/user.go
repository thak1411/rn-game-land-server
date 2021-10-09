package usecase

import (
	"github.com/google/uuid"
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/util"
)

type UserUsecase interface {
	CreateUser(model.User) error
	UpdateUser(model.User) error
	DeleteUser(int) error
	GetAllUser() ([]model.User, error)
	CheckUser(string, string) (bool, error)
}

type UserUC struct {
	db database.UserDatabase
}

func (uc *UserUC) CreateUser(user model.User) error {
	// Inject Salt & Hasing Password //
	user.Salt = uuid.New().String()
	user.Password = util.Encrypt(user.Password, user.Salt)
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

func (uc *UserUC) CheckUser(username, password string) (bool, error) {
	user, err := uc.db.GetUser(username)
	if err != nil {
		return false, err
	}
	password = util.Encrypt(password, user.Salt)
	return user.Id != -1 && user.Password == password, nil
}

func NewUser(db database.UserDatabase) UserUsecase {
	return &UserUC{db}
}
