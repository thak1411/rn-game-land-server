package usecase

import (
	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
)

type HubUsecase interface {
	RunHub()
	GetChatHub() *model.ChatHub
}

type HubUC struct {
	Hub *model.ChatHub
}

func (uc *HubUC) RunHub() {
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

func (uc *HubUC) GetChatHub() *model.ChatHub {
	return uc.Hub
}

func NewHub() HubUsecase {
	hub := &model.ChatHub{
		Clients:    make(map[*model.ChatClient]bool),
		Register:   make(chan *model.ChatClient),
		UnRegister: make(chan *model.ChatClient),
		Broadcast:  make(chan []byte),
		LastLog:    make([]string, 0, config.LastLogLen),
	}
	return &HubUC{hub}
}
