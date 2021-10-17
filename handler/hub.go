package handler

import (
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
)

type HubHandler struct {
	cuc usecase.ChatHubUsecase
	nuc usecase.NoticeHubUsecase
}

func (h *HubHandler) RunHub() {
	go h.cuc.RunChatHub()
	go h.nuc.RunNoticeHub()
}

func (h *HubHandler) GetChatHub() *model.ChatHub {
	return h.cuc.GetChatHub()
}

func (h *HubHandler) GetNoticeHub() *model.NoticeHub {
	return h.nuc.GetNoticeHub()
}

func NewHub(cuc usecase.ChatHubUsecase, nuc usecase.NoticeHubUsecase) *HubHandler {
	return &HubHandler{cuc, nuc}
}
