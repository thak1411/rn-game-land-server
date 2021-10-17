package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thak1411/rn-game-land-server/usecase"
)

type GameHandler struct {
	uc usecase.GameUsecase
}

func (h *GameHandler) GetGamelist(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		gamelist, err := h.uc.GetGameList()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(gamelist); err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func NewGame(uc usecase.GameUsecase) *GameHandler {
	return &GameHandler{uc}
}
