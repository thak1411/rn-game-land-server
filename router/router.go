package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/middleware"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func New() *http.ServeMux {
	hubUsecase := usecase.NewHub()
	hub := handler.NewHub(hubUsecase)
	hub.RunHub()

	clientUsecase := usecase.NewClient()
	client := handler.NewClient(clientUsecase)

	mux := http.NewServeMux()
	mux.HandleFunc("/ws/chat/connect", middleware.TokenParse(func(w http.ResponseWriter, r *http.Request) {
		client.WSChatServe(hub.GetChatHub(), w, r)
	}))

	userRouter := NewUser()
	gameRouter := NewGame()

	mux.Handle("/api/user/", http.StripPrefix("/api/user", userRouter))
	mux.Handle("/api/game/", http.StripPrefix("/api/game", gameRouter))
	return mux
}
