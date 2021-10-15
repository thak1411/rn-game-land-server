package main

import (
	"log"
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/middleware"
	"github.com/thak1411/rn-game-land-server/router"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func main() {
	hubUsecase := usecase.NewHub()
	hub := handler.NewHub(hubUsecase)
	hub.RunHub()

	clientUsecase := usecase.NewClient()
	client := handler.NewClient(clientUsecase)

	log.Println("Server Start! Listening Port", config.Port)
	router := router.New()
	router.HandleFunc("/ws/chat/connect", middleware.TokenParse(func(w http.ResponseWriter, r *http.Request) {
		client.WSChatServe(hub.GetChatHub(), w, r)
	}))
	http.ListenAndServe(config.Port, router)
}
