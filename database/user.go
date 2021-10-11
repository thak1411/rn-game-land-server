package database

import (
	"errors"
	"fmt"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
)

type UserDatabase interface {
	Create(model.User) error
	Update(model.User) error
	Delete(int) error
	GetAll() ([]model.User, error)
	GetUser(string) (model.User, error)
	GetUserId(string) (int, error)
	GetUserById(int) (model.User, error)
	AddFriend(int, int) error
}

type UserDB struct {
	users  map[int]model.User
	nextID int
}

func (db *UserDB) Create(user model.User) error {
	find, err := db.GetUser(user.Username)
	if err != nil {
		return err
	}
	if find.Id != -1 {
		return errors.New("duplicated username")
	}
	user.Id = db.nextID
	db.users[user.Id] = user
	db.nextID++
	return nil
}

func (db *UserDB) Update(user model.User) error {
	if _, ok := db.users[user.Id]; !ok {
		return errors.New("[UPDATE] There is No User With ID - " + fmt.Sprint(user.Id))
	}
	db.users[user.Id] = user
	return nil
}

func (db *UserDB) Delete(id int) error {
	if _, ok := db.users[id]; !ok {
		return errors.New("[UPDATE] There is No User With ID - " + fmt.Sprint(id))
	}
	delete(db.users, id)
	return nil
}

func (db *UserDB) GetAll() ([]model.User, error) {
	res := make([]model.User, 0, len(db.users))
	for _, v := range db.users {
		res = append(res, v)
	}
	return res, nil
}

func (db *UserDB) GetUser(username string) (model.User, error) {
	for _, v := range db.users {
		if v.Username == username {
			return v, nil
		}
	}
	return model.User{Id: -1}, nil
}

func (db *UserDB) GetUserId(username string) (int, error) {
	for _, v := range db.users {
		if v.Username == username {
			return v.Id, nil
		}
	}
	return -1, nil
}

func (db *UserDB) GetUserById(id int) (model.User, error) {
	user, ok := db.users[id]
	if ok {
		return user, nil
	}
	return model.User{Id: -1}, nil
}

func (db *UserDB) AddFriend(userId, targetId int) error {
	user, ok := db.users[userId]
	if !ok {
		return errors.New("user not found")
	}
	_, ok = db.users[targetId]
	if !ok {
		return errors.New("target not found")
	}
	user.Friend = append(user.Friend, targetId)
	db.users[userId] = user
	return nil
}

func NewUser() UserDatabase {
	return &UserDB{
		users: map[int]model.User{
			0: {Id: 0, Role: config.RoleAdmin, Name: "admin", Username: "admin", Salt: "admin_salt", Password: "892738161086b314334f88d661aa6e7bab7c825c34bf55222811dad46cdbf724"}, // pass: admin
		},
		nextID: 1,
	}
}
