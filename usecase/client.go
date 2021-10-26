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
	"github.com/thak1411/rn-game-land-server/games"
	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/util"
)

var badRequest = []byte(`{"code":400, "message": "bad request"}`)
var internalError = []byte(`{"code":500, "message": "internal server error"}`)
var unauthorizedError = []byte(`{"code":401, "message": "unauthorized behavior"}`)

type ClientUsecase interface {
	ClientReader(*model.WsClient)
	ClientWriter(*model.WsClient)
}

type ClientUC struct {
	gamedb      memorydb.GameDatabase
	userdb      database.UserDatabase
	gameHanlder games.GameHandler
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
		RoomName   string `json:"roomName"`
		GameName   string `json:"gameName"`
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
	res.RoomName = room.Name
	res.GameName = room.GameName
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
		// log.Printf("leave error")
		// client.Send <- internalError
		return
	}
}

type RejectInviteMessage struct {
	Code    int `json:"code"`
	Message struct {
		RoomId int `json:"roomId"`
	} `json:"message"`
}

type RejectInviteResponse struct {
	Code    int `json:"code"`
	Message struct {
		UserId int `json:"userId"`
		RoomId int `json:"roomId"`
	} `json:"message"`
}

func SendRejectInvite(uc *ClientUC, client *model.WsClient, message *RejectInviteMessage) {
	response := &RejectInviteResponse{}

	response.Code = 204
	res := &response.Message
	msg := &message.Message

	room, err := uc.gamedb.GetRoom(msg.RoomId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	cnt, err := uc.gamedb.DeleteInviteMessage(client.Id, msg.RoomId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	if cnt == 0 {
		log.Printf("not found to rejecting invite")
		client.Send <- badRequest
		return
	}

	del, err := uc.gamedb.DeleteRoomPlayer(msg.RoomId, client.Id)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	var targetsId []int
	if del {
		for _, v := range room.Player {
			if v.IsOnline {
				targetsId = append(targetsId, v.Id)
			}
		}
	}

	res.UserId = client.Id
	res.RoomId = msg.RoomId

	narrowMsg, err := json.Marshal(response)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	// x - TODO: must be spliting reject & remove protocol //
	narrowHandler := &model.NarrowcastHandler{
		Response: narrowMsg,
		Targets:  targetsId,
	}
	client.Hub.Narrowcast <- narrowHandler

	response.Code = 205
	narrowMsg2, err := json.Marshal(response)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	narrowHandler2 := &model.NarrowcastHandler{
		Response: narrowMsg2,
		Targets:  []int{client.Id},
	}
	client.Hub.Narrowcast <- narrowHandler2
}

type StartMessage struct {
	Code    int `json:"code"`
	Message struct {
		RoomId int `json:"roomId"`
	} `json:"message"`
}

type StartResponse struct {
	Code    int `json:"code"`
	Message struct {
		Room *model.Room `json:"room"`
	} `json:"message"`
}

func SendStart(uc *ClientUC, client *model.WsClient, message *StartMessage) {
	response := &StartResponse{}

	response.Code = 250
	res := &response.Message
	msg := &message.Message

	// get room //
	room, err := uc.gamedb.GetRoom(msg.RoomId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}
	// authorization //
	if room.Owner != client.Id {
		log.Printf("only owner can start game")
		client.Send <- unauthorizedError
		return
	}

	var targetsId []int
	for _, v := range room.Player {
		if v.IsOnline {
			targetsId = append(targetsId, v.Id)
		}
	}
	// delete offline player //
	for _, v := range room.Player {
		if v.IsOnline {
			continue
		}
		delResponse := &RejectInviteResponse{}

		delResponse.Code = 204
		delRes := &delResponse.Message
		_, err := uc.gamedb.DeleteInviteMessage(v.Id, msg.RoomId)
		if err != nil {
			log.Printf("error: %v", err)
			client.Send <- internalError
			return
		}

		del, err := uc.gamedb.DeleteRoomPlayer(msg.RoomId, v.Id)
		if err != nil {
			log.Printf("error: %v", err)
			client.Send <- internalError
			return
		}

		if del {
			delRes.UserId = v.Id
			delRes.RoomId = msg.RoomId

			narrowMsg, err := json.Marshal(delResponse)
			if err != nil {
				log.Printf("error: %v", err)
				client.Send <- internalError
				return
			}
			// x - TODO: must be spliting reject & remove protocol //
			narrowHandler := &model.NarrowcastHandler{
				Response: narrowMsg,
				// Targets:  append(targetsId, v.Id), // change this //
				Targets: targetsId, // change this //
			}
			client.Hub.Narrowcast <- narrowHandler

			delResponse.Code = 205
			narrowMsg2, err := json.Marshal(delResponse)
			if err != nil {
				log.Printf("error: %v", err)
				client.Send <- internalError
				return
			}
			narrowHandler2 := &model.NarrowcastHandler{
				Response: narrowMsg2,
				Targets:  []int{v.Id},
			}
			client.Hub.Narrowcast <- narrowHandler2
		}
	}

	// shuffle player & start //
	uc.gamedb.ShufflePlayer(msg.RoomId)

	_, err = uc.gamedb.SetGameStart(msg.RoomId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	room, err = uc.gamedb.GetRoom(msg.RoomId)
	if err != nil {
		log.Printf("error: %v", err)
		client.Send <- internalError
		return
	}

	res.Room = room

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

	uc.gameHanlder.Run(client.Hub, room)
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
	case model.REQ_REJECT:
		message := &RejectInviteMessage{}
		if err := util.BindJson(msg, message); err != nil {
			log.Printf("error: %v", err)
			client.Send <- badRequest
			return
		}
		SendRejectInvite(uc, client, message)
	case model.REQ_START:
		message := &StartMessage{}
		if err := util.BindJson(msg, message); err != nil {
			log.Printf("error: %v", err)
			client.Send <- badRequest
			return
		}
		SendStart(uc, client, message)
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

func NewClient(gamedb memorydb.GameDatabase, userdb database.UserDatabase, gameHandler games.GameHandler) ClientUsecase {
	return &ClientUC{gamedb, userdb, gameHandler}
}
