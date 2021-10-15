package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/util"
)

type ClientUsecase interface {
	ChatClientReader(*model.ChatClient)
	ChatClientWriter(*model.ChatClient)
}

type ClientUC struct{}

type ChatResponse struct {
	Code    int `json:"code"`
	Message struct {
		Id       int    `json:"id"`
		Time     string `json:"time"`
		Message  string `json:"message"`
		Username string `json:"username"`
	} `json:"message"`
}

func ChatHandler(uc *ClientUC, client *model.ChatClient, message *model.WsDefaultMessage) {
	response := &ChatResponse{Code: 200}
	errorRes := `{"code":500,"message":"Internal Server Error"}`
	switch message.Code {
	case 90: // Warnings: Change Code to Const Var //
		t := time.Now()
		response.Message.Id = client.Id
		response.Message.Time = fmt.Sprint(t.Hour()) + ":" + fmt.Sprint(t.Minute()) + ":" + fmt.Sprint(t.Second())
		response.Message.Message = message.Message
		response.Message.Username = client.Username

		msg, err := json.Marshal(response)
		if err != nil {
			client.Send <- []byte(errorRes)
			return
		}
		client.Hub.Broadcast <- msg
	}
}

func (uc *ClientUC) ChatClientReader(client *model.ChatClient) {
	defer func() {
		client.Hub.UnRegister <- client
		client.Conn.Close()
	}()
	client.Conn.SetReadLimit(config.MaxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(config.PongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })
	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, []byte("\n"), []byte(" "), -1))
		message := &model.WsDefaultMessage{}
		if err := util.BindJson(msg, message); err != nil {
			log.Printf("error: %v", err)
			break
		}
		ChatHandler(uc, client, message)
	}
}

func (uc *ClientUC) ChatClientWriter(client *model.ChatClient) {
	ticker := time.NewTicker(config.PingPeriod)
	defer func() {
		ticker.Stop()
		client.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			client.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
			if !ok { // hub close channel
				client.Conn.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			w, err := client.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// TODO: Update Send Data All in One //

			// n := len(client.Send)
			// for i := 0; i < n; i++ {
			// 	w.Write([]byte("\n"))
			// 	w.Write(<-client.Send)
			// }
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func NewClient() ClientUsecase {
	return &ClientUC{}
}
