package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/middleware"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func NewGame(gameDatabase memorydb.GameDatabase) *http.ServeMux {
	// authAdmin := middleware.AuthAdmin
	tokenDecode := middleware.TokenDecode

	mux := http.NewServeMux()

	// gameDatabase := memorydb.NewGame()
	gameUsecase := usecase.NewGame(gameDatabase)
	gameHandler := handler.NewGame(gameUsecase)

	mux.HandleFunc("/gamelist", gameHandler.GetGamelist)
	mux.HandleFunc("/create-room", tokenDecode(gameHandler.CreateRoom))
	mux.HandleFunc("/room", tokenDecode(gameHandler.GetRoom))
	return mux
}
