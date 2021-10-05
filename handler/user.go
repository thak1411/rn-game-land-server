package handler

import (
	"fmt"
	"net/http"

	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
)

type UserHandler struct {
	uc usecase.UserUsecase
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.uc.CreateUser(model.User{})
	fmt.Fprint(w, "API ROUTER TEST")
}

func NewUser(uc usecase.UserUsecase) *UserHandler {
	return &UserHandler{uc}
}
