package usecase

import (
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
	GetUserId(string) (int, error)
	GetUserIdByName(string) (int, error)
	GetUser(string) (model.User, error)
	GetUserById(int) (model.User, error)
	AddFriend(int, int) error
	RemoveFriend(int, int) error
	IsMyFriend(int, int) (bool, error)
	GetFriend(int) ([]model.User, error)
}

type UserUC struct {
	db database.UserDatabase
}

func (uc *UserUC) CreateUser(user model.User) error {
	user.Salt = util.NewUuid()
	user.Password = util.Encrypt(user.Password, user.Salt)
	return uc.db.Create(user)
}

func (uc *UserUC) UpdateUser(user model.User) error {
	user.Password = util.Encrypt(user.Password, user.Salt)
	return uc.db.Update(user)
}

func (uc *UserUC) DeleteUser(id int) error {
	return uc.db.Delete(id)
}

func (uc *UserUC) GetAllUser() ([]model.User, error) {
	return uc.db.GetAll()
}

func (uc *UserUC) GetUserId(username string) (int, error) {
	return uc.db.GetUserId(username)
}

func (uc *UserUC) GetUserIdByName(name string) (int, error) {
	return uc.db.GetUserIdByName(name)
}

func (uc *UserUC) GetUser(username string) (model.User, error) {
	return uc.db.GetUser(username)
}

func (uc *UserUC) CheckUser(username, password string) (bool, error) {
	user, err := uc.db.GetUser(username)
	if err != nil {
		return false, err
	}
	password = util.Encrypt(password, user.Salt)
	return user.Id != -1 && user.Password == password, nil
}

func (uc *UserUC) GetUserById(id int) (model.User, error) {
	return uc.db.GetUserById(id)
}

func (uc *UserUC) AddFriend(userId, targetId int) error {
	return uc.db.AddFriend(userId, targetId)
}

func (uc *UserUC) RemoveFriend(userId, targetId int) error {
	return uc.db.RemoveFriend(userId, targetId)
}

func (uc *UserUC) IsMyFriend(userId, targetId int) (bool, error) {
	return uc.db.IsMyFriend(userId, targetId)
}

func (uc *UserUC) GetFriend(userId int) ([]model.User, error) {
	return uc.db.GetFriend(userId)
}

func NewUser(db database.UserDatabase) UserUsecase {
	return &UserUC{db}
}
