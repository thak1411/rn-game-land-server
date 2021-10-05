package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/handler"
)

func NewIndex() *http.ServeMux {
	mux := http.NewServeMux()

	indexHandler := handler.NewIndex()

	mux.HandleFunc("/", indexHandler.ServeHTML)
	return mux
}
