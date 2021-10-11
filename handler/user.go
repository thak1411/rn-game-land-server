package handler

import (
	"fmt"
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
	"github.com/thak1411/rn-game-land-server/util"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

type UserForm struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body UserForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		user := model.User{
			Role:     config.RoleBasic,
			Name:     body.Name,
			Username: body.Username,
			Password: body.Password,
		}
		if err := h.uc.CreateUser(user); err != nil {
			fmt.Fprint(w, "Duplicated Username")
			return
		}
		fmt.Fprint(w, "Success Create User")
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := h.uc.GetAllUser()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "User List: %+v", users)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body LoginForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		flag, err := h.uc.CheckUser(body.Username, body.Password)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		if flag {
			user, err := h.uc.GetUser(body.Username)
			if err != nil {
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}
			token, err := util.CreateToken(user)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			cookie := &http.Cookie{
				Name:     config.Session,
				Value:    token,
				HttpOnly: true,
			}
			http.SetCookie(w, cookie)
			w.WriteHeader(200)
			return
		}
		w.WriteHeader(401)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

type FriendForm struct {
	Target   string `json:"target"`
	Username string `json:"username"`
}

func (h *UserHandler) AddFriend(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body FriendForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)
		if token.Username != body.Username {
			http.Error(w, "unahthorized token", http.StatusUnauthorized)
			return
		}
		target, err := h.uc.GetUser(body.Target)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		fmt.Println(token.Id, "->", target.Id)
		if err := h.uc.AddFriend(token.Id, target.Id); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NewUser(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc}
}
