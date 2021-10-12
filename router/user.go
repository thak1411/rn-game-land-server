package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/middleware"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func NewUser() *http.ServeMux {
	authAdmin := middleware.AuthAdmin
	tokenDecode := middleware.TokenDecode

	mux := http.NewServeMux()

	userDatabase := database.NewUser()
	userUsecase := usecase.NewUser(userDatabase)
	userHandler := handler.NewUser(userUsecase)

	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/logout", tokenDecode(userHandler.Logout))
	mux.HandleFunc("/user", userHandler.CreateUser)
	mux.HandleFunc("/add-friend", tokenDecode(userHandler.AddFriend))
	mux.HandleFunc("/all-user", tokenDecode(authAdmin(userHandler.GetAllUser)))
	mux.HandleFunc("/profile", tokenDecode(userHandler.GetUserProfile))
	return mux
}
