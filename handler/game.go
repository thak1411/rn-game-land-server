package handler

import (
	"encoding/json"
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/usecase"
	"github.com/thak1411/rn-game-land-server/util"
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

type RoomForm struct {
	GameId int    `json:"gameId"`
	Option string `json:"option"`
}

func (h *GameHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body RoomForm
		if err := util.BindBody(r.Body, &body); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		iToken := r.Context().Value(config.Session)
		token := iToken.(model.AuthTokenClaims)

		room, err := h.uc.CreateRoom(token.Id, body.GameId, token.Name+"'s Room", body.Option, token.Name)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(room); err != nil {
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
