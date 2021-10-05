package main

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/router"
)

func main() {
	router := router.New()
	http.ListenAndServe(config.Port, router)
}
