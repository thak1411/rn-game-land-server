package usecase

import (
	"github.com/thak1411/rn-game-land-server/model"
)

// type ChatHubUsecase interface {
// 	RunChatHub()
// 	GetChatHub() *model.ChatHub
// }

// type ChatHubUC struct {
// 	Hub *model.ChatHub
// }

// func (uc *ChatHubUC) RunChatHub() {
// 	for {
// 		select {
// 		case client := <-uc.Hub.Register:
// 			uc.Hub.Clients[client] = true
// 			for _, v := range uc.Hub.LastLog {
// 				client.Send <- []byte(v)
// 			}
// 		case client := <-uc.Hub.UnRegister:
// 			if _, ok := uc.Hub.Clients[client]; ok {
// 				delete(uc.Hub.Clients, client)
// 				close(client.Send)
// 			}
// 		case message := <-uc.Hub.Broadcast:
// 			for v := range uc.Hub.Clients {
// 				v.Send <- message
// 			}
// 			if len(uc.Hub.LastLog) < config.LastLogLen {
// 				uc.Hub.LastLog = append(uc.Hub.LastLog, string(message))
// 			} else {
// 				uc.Hub.LastLog = append(uc.Hub.LastLog[1:], string(message))
// 			}
// 		}
// 	}
// }

// func (uc *ChatHubUC) GetChatHub() *model.ChatHub {
// 	return uc.Hub
// }

// func NewChatHub() ChatHubUsecase {
// 	hub := &model.ChatHub{
// 		Clients:    make(map[*model.ChatClient]bool),
// 		Register:   make(chan *model.ChatClient),
// 		UnRegister: make(chan *model.ChatClient),
// 		Broadcast:  make(chan []byte),
// 		LastLog:    make([]string, 0, config.LastLogLen),
// 	}
// 	return &ChatHubUC{hub}
// }

// type NoticeHubUsecase interface {
// 	RunNoticeHub()
// 	GetNoticeHub() *model.NoticeHub
// }

// type NoticeHubUC struct {
// 	Hub *model.NoticeHub
// }

// type RetNoticeMessage struct {
// 	Code    int `json:"code"`
// 	Message struct {
// 		From   int `json:"from"`
// 		RoomId int `json:"roomId"`
// 	} `json:"message"`
// }

// func InviteToString(invite *model.InviteForm) []byte {
// 	ret := &RetNoticeMessage{Code: 200}
// 	emsg := `{"code":200,"message":"internal server error"}`
// 	ret.Code = 200
// 	ret.Message.From = invite.From
// 	ret.Message.RoomId = invite.RoomId
// 	b, err := json.Marshal(ret)
// 	if err != nil {
// 		return []byte(emsg)
// 	}
// 	return b
// }

// type RetInviteMessage struct {
// 	Code    int `json:"code"`
// 	Message struct {
// 		UserId   int    `json:"userId"`
// 		UserName string `json:"userName"`
// 	} `json:"message"`
// }

// func InviteToString2(invite *model.InviteForm) []byte {
// 	ret := &RetInviteMessage{Code: 203}
// 	emsg := `{"code":203,"message":"internal server error"}`
// 	ret.Code = 203
// 	ret.Message.UserId = invite.TargetId
// 	ret.Message.UserName = invite.TargetName
// 	b, err := json.Marshal(ret)
// 	if err != nil {
// 		return []byte(emsg)
// 	}
// 	return b
// }

// type RetJoinMessage struct {
// 	Code    int `json:"code"`
// 	Message struct {
// 		UserId int `json:"userId"`
// 	} `json:"message"`
// }

// func JoinToString(userId int) []byte {
// 	ret := &RetJoinMessage{Code: 201}
// 	emsg := `{"code":201,"message":"internal server error"}`
// 	ret.Code = 201
// 	ret.Message.UserId = userId
// 	b, err := json.Marshal(ret)
// 	if err != nil {
// 		return []byte(emsg)
// 	}
// 	return b
// }

// type RetLeaveMessage struct {
// 	Code    int `json:"code"`
// 	Message struct {
// 		UserId int `json:"userId"`
// 	} `json:"message"`
// }

// func LeaveToString(userId int) []byte {
// 	ret := &RetLeaveMessage{Code: 202}
// 	emsg := `{"code":202,"message":"internal server error"}`
// 	ret.Code = 202
// 	ret.Message.UserId = userId
// 	b, err := json.Marshal(ret)
// 	if err != nil {
// 		return []byte(emsg)
// 	}
// 	return b
// }

// func (uc *NoticeHubUC) RunNoticeHub() {
// 	for {
// 		select {
// 		case client := <-uc.Hub.Register:
// 			uc.Hub.Clients[client.Id] = client
// 			for _, v := range uc.Hub.InviteLog[client.Id] {
// 				uc.Hub.Clients[client.Id].Send <- InviteToString(v)
// 			}
// 		case client := <-uc.Hub.UnRegister:
// 			if _, ok := uc.Hub.Clients[client.Id]; ok {
// 				delete(uc.Hub.Clients, client.Id)
// 				close(client.Send)
// 			}
// 		case msg := <-uc.Hub.Invite:
// 			if _, ok := uc.Hub.Clients[msg.TargetId]; ok {
// 				uc.Hub.Clients[msg.TargetId].Send <- InviteToString(msg)
// 			}
// 			for _, v := range msg.TargetsId {
// 				if _, ok := uc.Hub.Clients[v]; ok {
// 					if uc.Hub.Clients[v].RoomId == msg.RoomId {
// 						uc.Hub.Clients[v].Send <- InviteToString2(msg)
// 					}
// 				}
// 			}
// 			uc.Hub.InviteLog[msg.TargetId] = append(uc.Hub.InviteLog[msg.TargetId], msg)
// 		case msg := <-uc.Hub.Join:
// 			if _, ok := uc.Hub.Clients[msg.UserId]; ok {
// 				uc.Hub.Clients[msg.UserId].RoomId = msg.RoomId
// 			}
// 			for _, v := range msg.TargetsId {
// 				if _, ok := uc.Hub.Clients[v]; ok {
// 					if uc.Hub.Clients[v].RoomId == msg.RoomId {
// 						uc.Hub.Clients[v].Send <- JoinToString(msg.UserId)
// 					}
// 				}
// 			}
// 		case msg := <-uc.Hub.Leave:
// 			if _, ok := uc.Hub.Clients[msg.UserId]; ok {
// 				uc.Hub.Clients[msg.UserId].RoomId = -1
// 			}
// 			for _, v := range msg.TargetsId {
// 				if _, ok := uc.Hub.Clients[v]; ok {
// 					if uc.Hub.Clients[v].RoomId == msg.RoomId {
// 						uc.Hub.Clients[v].Send <- LeaveToString(msg.UserId)
// 					}
// 				}
// 			}
// 		}
// 	}
// }

// func (uc *NoticeHubUC) GetNoticeHub() *model.NoticeHub {
// 	return uc.Hub
// }

// func NewNoticeHub() NoticeHubUsecase {
// 	hub := &model.NoticeHub{
// 		Clients:    make(map[int]*model.NoticeClient),
// 		Register:   make(chan *model.NoticeClient),
// 		UnRegister: make(chan *model.NoticeClient),
// 		Invite:     make(chan *model.InviteForm),
// 		InviteLog:  make(map[int][]*model.InviteForm),
// 		Join:       make(chan *model.JoinForm),
// 		Leave:      make(chan *model.LeaveForm),
// 	}
// 	return &NoticeHubUC{hub}
// }

type HubUsecase interface {
	RunHub()
	GetHub() *model.WsHub
}

type HubUC struct {
	Hub *model.WsHub
}

func (uc *HubUC) GetHub() *model.WsHub {
	return uc.Hub
}

func (uc *HubUC) RunHub() {
	for {
		select {
		case client := <-uc.Hub.Register:
			uc.Hub.Clients[client.Id] = client
			for _, v := range uc.Hub.BroadcastLog {
				client.Send <- v
			}
		case client := <-uc.Hub.UnRegister:
			if _, ok := uc.Hub.Clients[client.Id]; ok {
				delete(uc.Hub.Clients, client.Id)
				close(client.Send)
			}
		case msg := <-uc.Hub.Broadcast:
			for _, v := range uc.Hub.Clients {
				v.Send <- msg
			}
			if len(uc.Hub.BroadcastLog) < model.BC_LOG_LEN {
				uc.Hub.BroadcastLog = append(uc.Hub.BroadcastLog, msg)
			} else {
				uc.Hub.BroadcastLog = append(uc.Hub.BroadcastLog[1:], msg)
			}
		case msg := <-uc.Hub.Narrowcast:
			for _, v := range msg.Targets {
				if _, ok := uc.Hub.Clients[v]; ok {
					uc.Hub.Clients[v].Send <- msg.Response
				}
			}
		}
	}
}

func NewHub() HubUsecase {
	hub := &model.WsHub{
		Clients:      make(map[int]*model.WsClient),
		Register:     make(chan *model.WsClient),
		UnRegister:   make(chan *model.WsClient),
		Broadcast:    make(chan []byte),
		Narrowcast:   make(chan *model.NarrowcastHandler),
		BroadcastLog: [][]byte{},
	}
	return &HubUC{hub}
}
