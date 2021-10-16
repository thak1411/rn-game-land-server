package main

import (
	"log"
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/router"
)

func main() {
	log.Println("Server Start! Listening Port", config.Port)
	router := router.New()

	http.ListenAndServe(config.Port, router)
}
