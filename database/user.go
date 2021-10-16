package database

import (
	"bufio"
	"errors"
	"fmt"
	"os"

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
	GetUserIdByName(string) (int, error)
	GetUserById(int) (model.User, error)
	AddFriend(int, int) error
	RemoveFriend(int, int) error
	IsMyFriend(int, int) (bool, error)
}

type UserDB struct {
	users  map[int]model.User
	nextID int
}

func (db *UserDB) Create(user model.User) error {
	if id, err := db.GetUserId(user.Username); err != nil || id != -1 {
		if err != nil {
			return err
		}
		return errors.New("duplicated username")
	}
	if id, err := db.GetUserIdByName(user.Name); err != nil || id != -1 {
		if err != nil {
			return err
		}
		return errors.New("duplicated name")
	}
	user.Id = db.nextID
	user.Friend = make(map[int]bool)
	db.users[user.Id] = user
	db.nextID++
	db.BackupDB()
	return nil
}

func (db *UserDB) Update(user model.User) error {
	if _, ok := db.users[user.Id]; !ok {
		return errors.New("[UPDATE] There is No User With ID - " + fmt.Sprint(user.Id))
	}
	db.users[user.Id] = user
	db.BackupDB()
	return nil
}

func (db *UserDB) Delete(id int) error {
	if _, ok := db.users[id]; !ok {
		return errors.New("[UPDATE] There is No User With ID - " + fmt.Sprint(id))
	}
	delete(db.users, id)
	db.BackupDB()
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

func (db *UserDB) GetUserIdByName(name string) (int, error) {
	for _, v := range db.users {
		if v.Name == name {
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
	user.Friend[targetId] = true
	db.users[userId] = user
	db.BackupDB()
	return nil
}

func (db *UserDB) RemoveFriend(userId, targetId int) error {
	user, ok := db.users[userId]
	if !ok {
		return errors.New("user not found")
	}
	_, ok = db.users[targetId]
	if !ok {
		return errors.New("target not found")
	}
	delete(user.Friend, targetId)
	db.users[userId] = user
	db.BackupDB()
	return nil
}

func (db *UserDB) IsMyFriend(userId, targetId int) (bool, error) {
	user, ok := db.users[userId]
	if !ok {
		return false, errors.New("user not found")
	}
	_, ok = db.users[targetId]
	if !ok {
		return false, errors.New("target not found")
	}
	return user.Friend[targetId], nil
}

func (db *UserDB) BackupDB() {
	out, err := os.Create("temp.dat")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	w := bufio.NewWriter(out)
	defer w.Flush()
	fmt.Fprintf(w, "%d %d\n", len(db.users), db.nextID)
	for _, v := range db.users {
		fmt.Fprintf(w, "%d %s %s %s %s %s\n", v.Id, v.Role, v.Name, v.Username, v.Salt, v.Password)
		fmt.Fprintf(w, "%d", len(v.Friend))
		for t, aw := range v.Friend {
			var _ = aw
			fmt.Fprintf(w, " %d", t)
		}
		fmt.Fprintf(w, "\n")
	}
}

func LoadDB() (*UserDB, bool) {
	db := &UserDB{users: make(map[int]model.User)}
	in, err := os.Open("temp.dat")
	if err != nil {
		return nil, false
	}
	defer in.Close()

	r := bufio.NewReader(in)
	var n int
	fmt.Fscan(r, &n, &db.nextID)
	for i := 0; i < n; i++ {
		user := model.User{}
		fmt.Fscan(r, &user.Id, &user.Role, &user.Name, &user.Username, &user.Salt, &user.Password)
		var m int
		fmt.Fscan(r, &m)
		user.Friend = make(map[int]bool)
		for j := 0; j < m; j++ {
			var id int
			fmt.Fscan(r, &id)
			user.Friend[id] = true
		}
		db.users[i] = user
	}
	return db, true
}

func NewUser() UserDatabase {
	res, ok := LoadDB()
	if ok {
		return res
	}
	db := &UserDB{
		users: map[int]model.User{
			0: {
				Id:       0,
				Role:     config.RoleAdmin,
				Name:     "admin",
				Username: "admin",
				Salt:     "admin_salt",
				Password: "892738161086b314334f88d661aa6e7bab7c825c34bf55222811dad46cdbf724",
				Friend:   make(map[int]bool),
			}, // pass: admin
		},
		nextID: 1,
	}
	db.BackupDB()
	return db
}
