package main

import (
	"log"
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/router"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func main() {
	hubUsecase := usecase.NewHub()
	hub := handler.NewHub(hubUsecase)
	hub.RunHub()

	clientUsecase := usecase.NewClient()
	ws := handler.NewWs(clientUsecase)

	log.Println("Server Start! Listening Port", config.Port)
	router := router.New()
	router.HandleFunc("/ws/connect", func(w http.ResponseWriter, r *http.Request) {
		ws.WebSocketServe(hub.GetHub(), w, r)
	})
	http.ListenAndServe(config.Port, router)
}
