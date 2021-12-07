package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/games"
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

	gameHandler := games.New(gameDatabase, hub.GetHub())
	userDatabase := database.NewUser()

	clientUsecase := usecase.NewClient(gameDatabase, userDatabase, gameHandler)
	client := handler.NewClient(clientUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws/connect", middleware.TokenParse(func(w http.ResponseWriter, r *http.Request) {
		client.WsServe(hub.GetHub(), w, r)
	}))

	userRouter := NewUser(userDatabase)
	gameRouter := NewGame(gameDatabase)

	mux.Handle("/api/user/", http.StripPrefix("/api/user", userRouter))
	mux.Handle("/api/game/", http.StripPrefix("/api/game", gameRouter))
	return mux
}
