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

func (h *ClientHandler) WSChatServe(hub *model.ChatHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	iToken := r.Context().Value(config.Session)
	token := iToken.(model.AuthTokenClaims)

	client := &model.ChatClient{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 4096),
		ChatUser: model.ChatUser{
			Id:       token.Id,
			Name:     token.Name,
			Username: token.Username,
		},
	}
	client.Hub.Register <- client

	go h.uc.ChatClientReader(client)
	go h.uc.ChatClientWriter(client)
}

func (h *ClientHandler) WSNoticeServe(hub *model.NoticeHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	iToken := r.Context().Value(config.Session)
	token := iToken.(model.AuthTokenClaims)

	client := &model.NoticeClient{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 4096),
		NoticeUser: model.NoticeUser{
			Id:       token.Id,
			Name:     token.Name,
			RoomId:   -1,
			Username: token.Username,
		},
	}
	client.Hub.Register <- client

	go h.uc.NoticeClientReader(client)
	go h.uc.NoticeClientWriter(client)
}

func NewClient(uc usecase.ClientUsecase) *ClientHandler {
	return &ClientHandler{uc}
}
