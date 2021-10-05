package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	userUsecase := usecase.NewUser()
	userHandler := handler.NewUser(userUsecase)

	mux.HandleFunc("/user", userHandler.NewUser)
	return mux
}
