package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/middleware"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func New() *http.ServeMux {
	gameDatabase := memorydb.NewGame()

	hubUsecase := usecase.NewHub(gameDatabase)
	hub := handler.NewHub(hubUsecase)
	hub.RunHub()

	userDatabase := database.NewUser()
	clientUsecase := usecase.NewClient(gameDatabase, userDatabase)
	client := handler.NewClient(clientUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", middleware.TokenParse(func(w http.ResponseWriter, r *http.Request) {
		client.WsServe(hub.GetHub(), w, r)
	}))

	userRouter := NewUser()
	gameRouter := NewGame()

	mux.Handle("/api/user/", http.StripPrefix("/api/user", userRouter))
	mux.Handle("/api/game/", http.StripPrefix("/api/game", gameRouter))
	return mux
}
