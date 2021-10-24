package usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/util"
)

// type ClientUsecase interface {
// 	ChatClientReader(*model.ChatClient)
// 	ChatClientWriter(*model.ChatClient)
// 	NoticeClientReader(*model.NoticeClient)
// 	NoticeClientWriter(*model.NoticeClient)
// }

// type ClientUC struct {
// 	db     database.GameDatabase
// 	userdb database.UserDatabase
// }

// type ChatResponse struct {
// 	Code    int `json:"code"`
// 	Message struct {
// 		Id      int    `json:"id"`
// 		Time    string `json:"time"`
// 		Name    string `json:"name"`
// 		Message string `json:"message"`
// 	} `json:"message"`
// }

// func ChatHandler(uc *ClientUC, client *model.ChatClient, message *model.WsDefaultMessage) {
// 	response := &ChatResponse{Code: 200}
// 	errorRes := `{"code":500,"message":"Internal Server Error"}`
// 	switch message.Code {
// 	case 90: // Warnings: Change Code to Const Var //
// 		t := time.Now()
// 		response.Message.Id = client.Id
// 		response.Message.Time = fmt.Sprint(t.Hour()) + ":" + fmt.Sprint(t.Minute()) + ":" + fmt.Sprint(t.Second())
// 		response.Message.Name = client.Name
// 		response.Message.Message = message.Message

// 		msg, err := json.Marshal(response)
// 		if err != nil {
// 			client.Send <- []byte(errorRes)
// 			return
// 		}
// 		client.Hub.Broadcast <- msg
// 	}
// }

// func (uc *ClientUC) ChatClientReader(client *model.ChatClient) {
// 	defer func() {
// 		client.Hub.UnRegister <- client
// 		client.Conn.Close()
// 	}()
// 	client.Conn.SetReadLimit(config.MaxMessageSize)
// 	client.Conn.SetReadDeadline(time.Now().Add(config.PongWait))
// 	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })
// 	for {
// 		_, msg, err := client.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		msg = bytes.TrimSpace(bytes.Replace(msg, []byte("\n"), []byte(" "), -1))
// 		message := &model.WsDefaultMessage{}
// 		if err := util.BindJson(msg, message); err != nil {
// 			log.Printf("error: %v", err)
// 			break
// 		}
// 		ChatHandler(uc, client, message)
// 	}
// }

// func (uc *ClientUC) ChatClientWriter(client *model.ChatClient) {
// 	ticker := time.NewTicker(config.PingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		client.Conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-client.Send:
// 			client.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
// 			if !ok { // hub close channel
// 				client.Conn.WriteMessage(websocket.CloseMessage, nil)
// 				return
// 			}
// 			w, err := client.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			// TODO: Update Send Data All in One //

// 			// n := len(client.Send)
// 			// for i := 0; i < n; i++ {
// 			// 	w.Write([]byte("\n"))
// 			// 	w.Write(<-client.Send)
// 			// }
// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			client.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
// 			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

// func NoticeHandler(uc *ClientUC, client *model.NoticeClient, message *model.WsDefaultMessage) {
// 	switch message.Code {
// 	case 50: // Warnings: Change Code to Const Var //
// 		msg := &model.InviteForm{}
// 		if err := json.Unmarshal([]byte(message.Message), msg); err != nil {
// 			log.Printf("error: %v", err)
// 			return
// 		}
// 		msg.From = client.Id
// 		room, err := uc.db.GetRoom(msg.RoomId)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			return
// 		}

// 		if room.Owner != msg.From {
// 			log.Printf("only owner can invite people")
// 			return
// 		}

// 		targetName, err := uc.userdb.GetNameById(msg.TargetId)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			return
// 		}

// 		uc.db.AppendRoomPlayer(msg.RoomId, msg.TargetId, targetName)
// 		msg.From = client.Id
// 		msg.TargetName = targetName

// 		for _, v := range room.Player {
// 			msg.TargetsId = append(msg.TargetsId, v.Id)
// 		}
// 		client.Hub.Invite <- msg
// 		// TODO: limit room size //
// 	case 51:
// 		roomId, err := strconv.Atoi(message.Message)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			return
// 		}
// 		join, err := uc.db.SetUserOnline(roomId, client.Id)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			return
// 		}
// 		if join {
// 			room, err := uc.db.GetRoom(roomId)
// 			if err != nil {
// 				log.Printf("error: %v", err)
// 				return
// 			}
// 			msg := &model.JoinForm{}
// 			msg.UserId = client.Id
// 			msg.RoomId = roomId
// 			for _, v := range room.Player {
// 				msg.TargetsId = append(msg.TargetsId, v.Id)
// 			}
// 			client.Hub.Join <- msg
// 		}
// 	}
// }

// func (uc *ClientUC) NoticeClientReader(client *model.NoticeClient) {
// 	defer func() {
// 		client.Hub.UnRegister <- client
// 		client.Conn.Close()
// 	}()
// 	client.Conn.SetReadLimit(config.MaxMessageSize)
// 	client.Conn.SetReadDeadline(time.Now().Add(config.PongWait))
// 	client.Conn.SetPongHandler(func(string) error { client.Conn.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })
// 	for {
// 		_, msg, err := client.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		msg = bytes.TrimSpace(bytes.Replace(msg, []byte("\n"), []byte(" "), -1))
// 		message := &model.WsDefaultMessage{}
// 		if err := util.BindJson(msg, message); err != nil {
// 			log.Printf("error: %v", err)
// 			break
// 		}
// 		NoticeHandler(uc, client, message)
// 	}
// }

// func (uc *ClientUC) NoticeClientWriter(client *model.NoticeClient) {
// 	ticker := time.NewTicker(config.PingPeriod)
// 	defer func() {
// 		// in disconnected //
// 		leave, room, err := uc.db.SetUserOffline(client.Id)
// 		if err != nil {
// 			log.Printf("error: %v", err)
// 			return
// 		}
// 		if leave {
// 			msg := &model.LeaveForm{}
// 			msg.UserId = client.Id
// 			msg.RoomId = room.Id
// 			for _, v := range room.Player {
// 				if v.Id == client.Id {
// 					continue
// 				}
// 				msg.TargetsId = append(msg.TargetsId, v.Id)
// 			}
// 			client.Hub.Leave <- msg
// 		}

// 		ticker.Stop()
// 		client.Conn.Close()
// 	}()
// 	for {
// 		select {
// 		case message, ok := <-client.Send:
// 			client.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
// 			if !ok { // hub close channel
// 				client.Conn.WriteMessage(websocket.CloseMessage, nil)
// 				return
// 			}
// 			w, err := client.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write(message)

// 			// TODO: Update Send Data All in One //

// 			// n := len(client.Send)
// 			// for i := 0; i < n; i++ {
// 			// 	w.Write([]byte("\n"))
// 			// 	w.Write(<-client.Send)
// 			// }
// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			client.Conn.SetWriteDeadline(time.Now().Add(config.WriteWait))
// 			if err := client.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }

var badRequest = []byte(`{"code":400, "message": "bad request"}`)
var internalError = []byte(`{"code":500, "message": "internal server error"}`)
var unauthorizedError = []byte(`{"code":401, "message": "unauthorized behavior"}`)

type ClientUsecase interface {
	ClientReader(*model.WsClient)
	ClientWriter(*model.WsClient)
}

type ClientUC struct {
	gamedb memorydb.GameDatabase
	userdb database.UserDatabase
}

type ChatMessage struct {
	Code    int `json:"code"`
	Message struct {
		Data string `json:"data"`
	} `json:"message"`
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

func SendChat(uc *ClientUC, client *model.WsClient, message *ChatMessage) {
	response := &ChatResponse{}
	t := time.Now()
	response.Code = model.RES_BROADCAST
	response.Message.Id = client.Id
	response.Message.Time = fmt.Sprint(t.Hour()) + ":" + fmt.Sprint(t.Minute()) + ":" + fmt.Sprint(t.Second())
	response.Message.Name = client.Name
	response.Message.Message = message.Message.Data

	msg, err := json.Marshal(response)
	if err != nil {
		client.Send <- internalError
		return
	}
	client.Hub.Broadcast <- msg
}

type InviteMessage struct {
	Code    int `json:"code"`
	Message struct {
		RoomId   int `json:"roomId"`
		TargetId int `json:"targetId"`
	} `json:"message"`
}

type InviteResponse struct {
	Code    int `json:"code"`
	Message struct {
		From       int    `json:"from"`
		RoomId     int    `json:"roomId"`
		FromName   string `json:"fromName"`
		TargetId   int    `json:"targetId"`
		TargetName string `json:"targetName"`
	} `json:"message"`
}

func SendInvite(uc *ClientUC, client *model.WsClient, message *InviteMessage) {
	response := &InviteResponse{}

	msg := &message.Message
	res := &response.Message

	response.Code = 203
	res.From = client.Id
	res.RoomId = msg.RoomId
	res.FromName = client.Name
	res.TargetId = msg.TargetId

	room, err := uc.gamedb.GetRoom(res.RoomId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	if room.Owner != res.From {
		log.Printf("only owner can invite people")
		client.Send <- unauthorizedError
		return
	}

	targetName, err := uc.userdb.GetNameById(res.TargetId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	uc.gamedb.AppendRoomPlayer(res.RoomId, res.TargetId, targetName)
	res.From = client.Id
	res.TargetName = targetName

	var targetsId []int
	for _, v := range room.Player {
		if v.IsOnline {
			targetsId = append(targetsId, v.Id)
		}
	}

	narrowMsg, err := json.Marshal(response)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	narrowHandler := &model.NarrowcastHandler{
		Response: narrowMsg,
		Targets:  targetsId,
	}
	client.Hub.Narrowcast <- narrowHandler

	response.Code = 200
	narrowMsg2, err := json.Marshal(response)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	narrowHandler2 := &model.NarrowcastHandler{
		Response: narrowMsg2,
		Targets:  []int{msg.TargetId},
	}
	uc.gamedb.AppendInviteMessage(msg.TargetId, narrowMsg2)
	client.Hub.Narrowcast <- narrowHandler2
	// TODO: limit room size //
}

type JoinMessage struct {
	Code    int `json:"code"`
	Message struct {
		RoomId int `json:"roomId"`
	} `json:"message"`
}

type JoinResponse struct {
	Code    int `json:"code"`
	Message struct {
		UserId int `json:"userId"`
		RoomId int `json:"roomId"`
	} `json:"message"`
}

func SendJoin(uc *ClientUC, client *model.WsClient, message *JoinMessage) {
	response := &JoinResponse{}

	msg := &message.Message
	res := &response.Message

	response.Code = 201

	join, err := uc.gamedb.SetUserOnline(msg.RoomId, client.Id)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	if join {
		client.RoomId = msg.RoomId

		room, err := uc.gamedb.GetRoom(msg.RoomId)
		if err != nil {
			log.Printf("error: %v", err)
			client.Send <- internalError
			return
		}
		res.UserId = client.Id
		res.RoomId = msg.RoomId

		var targetsId []int
		for _, v := range room.Player {
			if v.IsOnline {
				targetsId = append(targetsId, v.Id)
			}
		}

		narrowMsg, err := json.Marshal(response)
		if err != nil {
			log.Printf("error: %v", err)
			client.Send <- internalError
			return
		}

		narrowHandler := &model.NarrowcastHandler{
			Response: narrowMsg,
			Targets:  targetsId,
		}
		_, err = uc.gamedb.DeleteInviteMessage(client.Id, msg.RoomId)
		if err != nil {
			log.Printf("error: %v", err)
		}
		client.Hub.Narrowcast <- narrowHandler
	} else {
		log.Printf("join error")
		client.Send <- internalError
		return
	}
}

type LeaveResponse struct {
	Code    int `json:"code"`
	Message struct {
		UserId int `json:"userId"`
		RoomId int `json:"roomId"`
	} `json:"message"`
}

func SendLeave(uc *ClientUC, client *model.WsClient) {
	// TODO: check user has room //
	if client.RoomId == -1 {
		return
	}

	response := &LeaveResponse{}

	response.Code = 202
	res := &response.Message

	leave, room, err := uc.gamedb.SetUserOffline(client.Id)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	if leave {
		res.UserId = client.Id
		res.RoomId = room.Id

		var targetsId []int
		for _, v := range room.Player {
			if v.Id != client.Id && v.IsOnline {
				targetsId = append(targetsId, v.Id)
			}
		}
		narrowMsg, err := json.Marshal(response)
		if err != nil {
			log.Printf("error: %v", err)
			client.Send <- internalError
			return
		}

		narrowHandler := &model.NarrowcastHandler{
			Response: narrowMsg,
			Targets:  targetsId,
		}
		client.Hub.Narrowcast <- narrowHandler
	} else {
		log.Printf("leave error")
		client.Send <- internalError
		return
	}
}

func WsHandler(uc *ClientUC, client *model.WsClient, msg []byte) {
	defaultMessage := &model.WsDefaultMessage{}
	if err := util.BindJson(msg, defaultMessage); err != nil {
		log.Printf("error: %v", err)
		client.Send <- badRequest
		return
	}

	switch defaultMessage.Code {
	case model.REQ_BROADCAST:
		message := &ChatMessage{}
		if err := util.BindJson(msg, message); err != nil {
			log.Printf("error: %v", err)
			client.Send <- badRequest
			return
		}
		SendChat(uc, client, message)
	case model.REQ_INVITE:
		message := &InviteMessage{}
		if err := util.BindJson(msg, message); err != nil {
			log.Printf("error: %v", err)
			client.Send <- badRequest
			return
		}
		SendInvite(uc, client, message)
	case model.REQ_JOIN:
		message := &JoinMessage{}
		if err := util.BindJson(msg, message); err != nil {
			log.Printf("error: %v", err)
			client.Send <- badRequest
			return
		}
		SendJoin(uc, client, message)
	}
}

func (uc *ClientUC) ClientReader(client *model.WsClient) {
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
		WsHandler(uc, client, msg)
	}
}

func (uc *ClientUC) ClientWriter(client *model.WsClient) {
	ticker := time.NewTicker(config.PingPeriod)
	defer func() {
		SendLeave(uc, client) // in disconnected send room leave event

		ticker.Stop()
		client.Conn.Close()
	}()

	// Send Saved Invite Message in Connection //
	inviteMsg, err := uc.gamedb.GetInviteMessage(client.Id)
	if err != nil {
		log.Printf("saved invite message error: %v", err)
		client.Send <- internalError
	} else {
		for _, v := range inviteMsg {
			client.Send <- v
		}
	}

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

func NewClient(gamedb memorydb.GameDatabase, userdb database.UserDatabase) ClientUsecase {
	return &ClientUC{gamedb, userdb}
}
