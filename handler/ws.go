package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
)

type WsHandler struct {
	uc usecase.ClientUsecase
}

func (h *WsHandler) WebSocketServe(hub *model.Hub, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &model.Client{Hub: hub, Conn: conn, Send: make(chan []byte, 512)}
	client.Hub.Register <- client

	go h.uc.ClientReader(client)
	go h.uc.ClientWriter(client)
}

func NewWs(uc usecase.ClientUsecase) *WsHandler {
	return &WsHandler{uc}
}
