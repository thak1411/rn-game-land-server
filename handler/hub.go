package handler

import (
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
)

type HubHandler struct {
	uc usecase.HubUsecase
}

func (h *HubHandler) RunHub() {
	go h.uc.RunHub()
}

func (h *HubHandler) GetHub() *model.Hub {
	return h.uc.GetHub()
}

func NewHub(uc usecase.HubUsecase) *HubHandler {
	return &HubHandler{uc}
}
