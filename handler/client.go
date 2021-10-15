package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8192,
	WriteBufferSize: 8192,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ClientHandler struct {
	uc usecase.ClientUsecase
}

func (h *ClientHandler) WSChatServe(hub *model.ChatHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &model.ChatClient{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 4096),
	}
	client.Hub.Register <- client

	go h.uc.ChatClientReader(client)
	go h.uc.ChatClientWriter(client)
}

func NewClient(uc usecase.ClientUsecase) *ClientHandler {
	return &ClientHandler{uc}
}