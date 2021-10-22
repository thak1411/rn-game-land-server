package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/config"
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

func (h *ClientHandler) WsServe(hub *model.WsHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	iToken := r.Context().Value(config.Session)
	token := iToken.(model.AuthTokenClaims)

	client := &model.WsClient{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 4096),
		WsUser: model.WsUser{
			Id:       token.Id,
			Name:     token.Name,
			RoomId:   -1,
			Username: token.Username,
		},
	}
	client.Hub.Register <- client

	go h.uc.ClientReader(client)
	go h.uc.ClientWriter(client)
}

func NewClient(uc usecase.ClientUsecase) *ClientHandler {
	return &ClientHandler{uc}
}
