package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func NewUser() *http.ServeMux {
	mux := http.NewServeMux()

	userDatabase := database.NewUser()
	userUsecase := usecase.NewUser(userDatabase)
	userHandler := handler.NewUser(userUsecase)

	mux.HandleFunc("/user", userHandler.CreateUser)
	mux.HandleFunc("/all-user", userHandler.GetAllUser)
	return mux
}
