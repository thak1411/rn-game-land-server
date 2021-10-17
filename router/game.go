package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/middleware"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func NewGame() *http.ServeMux {
	// authAdmin := middleware.AuthAdmin
	tokenDecode := middleware.TokenDecode

	mux := http.NewServeMux()

	gameDatabase := database.NewGame()
	gameUsecase := usecase.NewGame(gameDatabase)
	gameHandler := handler.NewGame(gameUsecase)

	mux.HandleFunc("/gamelist", gameHandler.GetGamelist)
	mux.HandleFunc("/create-room", tokenDecode(gameHandler.CreateRoom))
	return mux
}
