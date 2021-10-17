package usecase

import (
	"encoding/json"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
)

type ChatHubUsecase interface {
	RunChatHub()
	GetChatHub() *model.ChatHub
}

type ChatHubUC struct {
	Hub *model.ChatHub
}

func (uc *ChatHubUC) RunChatHub() {
	for {
		select {
		case client := <-uc.Hub.Register:
			uc.Hub.Clients[client] = true
			for _, v := range uc.Hub.LastLog {
				client.Send <- []byte(v)
			}
		case client := <-uc.Hub.UnRegister:
			if _, ok := uc.Hub.Clients[client]; ok {
				delete(uc.Hub.Clients, client)
				close(client.Send)
			}
		case message := <-uc.Hub.Broadcast:
			for v := range uc.Hub.Clients {
				v.Send <- message
			}
			if len(uc.Hub.LastLog) < config.LastLogLen {
				uc.Hub.LastLog = append(uc.Hub.LastLog, string(message))
			} else {
				uc.Hub.LastLog = append(uc.Hub.LastLog[1:], string(message))
			}
		}
	}
}

func (uc *ChatHubUC) GetChatHub() *model.ChatHub {
	return uc.Hub
}

func NewChatHub() ChatHubUsecase {
	hub := &model.ChatHub{
		Clients:    make(map[*model.ChatClient]bool),
		Register:   make(chan *model.ChatClient),
		UnRegister: make(chan *model.ChatClient),
		Broadcast:  make(chan []byte),
		LastLog:    make([]string, 0, config.LastLogLen),
	}
	return &ChatHubUC{hub}
}

type NoticeHubUsecase interface {
	RunNoticeHub()
	GetNoticeHub() *model.NoticeHub
}

type NoticeHubUC struct {
	Hub *model.NoticeHub
}

type RetNoticeMessage struct {
	Code    int `json:"code"`
	Message struct {
		From   int `json:"from"`
		RoomId int `json:"roomId"`
	} `json:"message"`
}

func InviteToString(invite *model.InviteForm) []byte {
	ret := &RetNoticeMessage{Code: 200}
	emsg := `{"code":200,"message":"internal server error"}`
	ret.Code = 200
	ret.Message.From = invite.From
	ret.Message.RoomId = invite.RoodId
	b, err := json.Marshal(ret)
	if err != nil {
		return []byte(emsg)
	}
	return b
}

func (uc *NoticeHubUC) RunNoticeHub() {
	for {
		select {
		case client := <-uc.Hub.Register:
			uc.Hub.Clients[client.Id] = client
			for _, v := range uc.Hub.InviteLog[client.Id] {
				uc.Hub.Clients[client.Id].Send <- InviteToString(v)
			}
		case client := <-uc.Hub.UnRegister:
			if _, ok := uc.Hub.Clients[client.Id]; ok {
				delete(uc.Hub.Clients, client.Id)
				close(client.Send)
			}
		case msg := <-uc.Hub.Invite:
			if _, ok := uc.Hub.Clients[msg.TargetId]; ok {
				uc.Hub.Clients[msg.TargetId].Send <- InviteToString(msg)
			}
			uc.Hub.InviteLog[msg.TargetId] = append(uc.Hub.InviteLog[msg.TargetId], msg)
		}
	}
}

func (uc *NoticeHubUC) GetNoticeHub() *model.NoticeHub {
	return uc.Hub
}

func NewNoticeHub() NoticeHubUsecase {
	hub := &model.NoticeHub{
		Clients:    make(map[int]*model.NoticeClient),
		Register:   make(chan *model.NoticeClient),
		UnRegister: make(chan *model.NoticeClient),
		Invite:     make(chan *model.InviteForm),
		InviteLog:  make(map[int][]*model.InviteForm),
	}
	return &NoticeHubUC{hub}
}
