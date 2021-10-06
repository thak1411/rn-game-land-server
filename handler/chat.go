package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/usecase"
)

type ChatHandler struct {
	uc usecase.ChatUsecase
}

func (h *ChatHandler) SocketTest(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer conn.Close()
	h.uc.SocketTest(conn)
}

func NewChat(uc usecase.ChatUsecase) *ChatHandler {
	return &ChatHandler{uc}
}
