package database

import (
	"errors"
	"fmt"

	"github.com/thak1411/rn-game-land-server/model"
)

type UserDatabase interface {
	Create(model.User) error
	Update(model.User) error
	Delete(int) error
	GetAll() ([]model.User, error)
}

type UserDB struct {
	users  map[int]model.User
	nextID int
}

func (h *UserDB) Create(user model.User) error {
	user.Id = h.nextID
	h.users[user.Id] = user
	h.nextID++
	return nil
}

func (h *UserDB) Update(user model.User) error {
	if _, ok := h.users[user.Id]; !ok {
		return errors.New("[UPDATE] There is No User With ID - " + fmt.Sprint(user.Id))
	}
	h.users[user.Id] = user
	return nil
}

func (h *UserDB) Delete(id int) error {
	if _, ok := h.users[id]; !ok {
		return errors.New("[UPDATE] There is No User With ID - " + fmt.Sprint(id))
	}
	delete(h.users, id)
	return nil
}

func (h *UserDB) GetAll() ([]model.User, error) {
	res := make([]model.User, 0, len(h.users))
	for _, v := range h.users {
		res = append(res, v)
	}
	return res, nil
}

func NewUser() UserDatabase {
	return &UserDB{
		users: map[int]model.User{
			0: {Id: 0, Name: "admin", Username: "admin", Password: "pass"},
		},
		nextID: 1,
	}
}
