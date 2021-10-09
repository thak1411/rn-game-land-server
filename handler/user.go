package handler

import (
	"fmt"
	"net/http"

	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
	"github.com/thak1411/rn-game-land-server/util"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		user := model.User{}
		err := util.Bind(r.Body, &user)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		h.uc.CreateUser(user)
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// h.uc.CheckUser()
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NewUser(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc}
}
