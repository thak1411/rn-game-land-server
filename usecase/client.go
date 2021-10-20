package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/util"
)

type ClientUsecase interface {
	ChatClientReader(*model.ChatClient)
	ChatClientWriter(*model.ChatClient)
	NoticeClientReader(*model.NoticeClient)
	NoticeClientWriter(*model.NoticeClient)
}

type ClientUC struct {
	db     database.GameDatabase
	userdb database.UserDatabase
}

type ChatResponse struct {
	Code    int `json:"code"`
	Message struct {
		Id      int    `json:"id"`
		Time    string `json:"time"`
		Name    string `json:"name"`
		Message string `json:"message"`
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
		response.Message.Name = client.Name
		response.Message.Message = message.Message

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

func NoticeHandler(uc *ClientUC, client *model.NoticeClient, message *model.WsDefaultMessage) {
	switch message.Code {
	case 50: // Warnings: Change Code to Const Var //
		msg := &model.InviteForm{}
		if err := json.Unmarshal([]byte(message.Message), msg); err != nil {
			log.Printf("error: %v", err)
			return
		}
		targetName, err := uc.userdb.GetNameById(msg.TargetId)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		uc.db.AppendRoomPlayer(msg.RoomId, msg.TargetId, targetName)
		msg.From = client.Id
		msg.TargetName = targetName

		room, err := uc.db.GetRoom(msg.RoomId)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		for _, v := range room.Player {
			msg.TargetsId = append(msg.TargetsId, v.Id)
		}
		client.Hub.Invite <- msg
		// TODO: limit room size //
	case 51:
		roomId, err := strconv.Atoi(message.Message)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		join, err := uc.db.SetUserOnline(roomId, client.Id)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		if join {
			room, err := uc.db.GetRoom(roomId)
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
			msg := &model.JoinForm{}
			msg.UserId = client.Id
			msg.RoomId = roomId
			for _, v := range room.Player {
				msg.TargetsId = append(msg.TargetsId, v.Id)
			}
			client.Hub.Join <- msg
		}
	}
}

func (uc *ClientUC) NoticeClientReader(client *model.NoticeClient) {
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
		NoticeHandler(uc, client, message)
	}
}

func (uc *ClientUC) NoticeClientWriter(client *model.NoticeClient) {
	ticker := time.NewTicker(config.PingPeriod)
	defer func() {
		// in disconnected //
		leave, room, err := uc.db.SetUserOffline(client.Id)
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		if leave {
			msg := &model.LeaveForm{}
			msg.UserId = client.Id
			msg.RoomId = room.Id
			for _, v := range room.Player {
				if v.Id == client.Id {
					continue
				}
				msg.TargetsId = append(msg.TargetsId, v.Id)
			}
			client.Hub.Leave <- msg
		}

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

func NewClient(db database.GameDatabase, userdb database.UserDatabase) ClientUsecase {
	return &ClientUC{db, userdb}
}
