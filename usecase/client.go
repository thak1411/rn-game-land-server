package usecase

import (
	"bytes"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
)

type ClientUsecase interface {
	ClientReader(*model.Client)
	ClientWriter(*model.Client)
}

type ClientUC struct{}

func (uc *ClientUC) ClientReader(client *model.Client) {
	defer func() {
		client.Hub.UnRegister <- client
		client.Conn.Close()
	}()
	client.Conn.SetReadLimit(config.MaxMessageSize)
	client.Conn.SetReadDeadline(time.Now().Add(config.PongWait))
	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))
		client.Send <- message
	}
}

func (uc *ClientUC) ClientWriter(client *model.Client) {
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
			w, err := client.Conn.NextWriter((websocket.TextMessage))
			if err != nil {
				return
			}
			w.Write(message)

			n := len(client.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-client.Send)
			}
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
