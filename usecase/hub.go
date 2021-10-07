package usecase

import "github.com/thak1411/rn-game-land-server/model"

type HubUsecase interface {
	RunHub()
	GetHub() *model.Hub
}

type HubUC struct {
	Hub *model.Hub
}

func (uc *HubUC) RunHub() {
	for {
		select {
		case client := <-uc.Hub.Register:
			uc.Hub.Clients[client] = "TEST_NAME"
		case client := <-uc.Hub.UnRegister:
			if _, ok := uc.Hub.Clients[client]; ok {
				delete(uc.Hub.Clients, client)
				close(client.Send)
			}
		}
	}
}

func (uc *HubUC) GetHub() *model.Hub {
	return uc.Hub
}

func NewHub() HubUsecase {
	hub := &model.Hub{
		Clients:    make(map[*model.Client]string),
		Register:   make(chan *model.Client),
		UnRegister: make(chan *model.Client),
	}
	return &HubUC{hub}
}
