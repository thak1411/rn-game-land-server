package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func NewUser() *http.ServeMux {
	mux := http.NewServeMux()

	userUsecase := usecase.NewUser()
	userHandler := handler.NewUser(userUsecase)

	mux.HandleFunc("/test", userHandler.CreateUser)
	return mux
}
