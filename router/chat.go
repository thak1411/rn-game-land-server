package router

import (
	"net/http"

	"github.com/thak1411/rn-game-land-server/handler"
	"github.com/thak1411/rn-game-land-server/usecase"
)

func NewChat() *http.ServeMux {
	mux := http.NewServeMux()

	chatUsecase := usecase.NewChat()
	chatHandler := handler.NewChat(chatUsecase)

	mux.HandleFunc("/test", chatHandler.SocketTest)
	return mux
}
