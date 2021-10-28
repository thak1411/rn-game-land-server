package games

import (
	"github.com/thak1411/rn-game-land-server/games/yahtzee"
	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/model"
)

type GameHandler interface {
	Run(*model.Room)
}

type Games struct {
	gamedb memorydb.GameDatabase
	hub    *model.WsHub
}

func (h *Games) Run(room *model.Room) {
	switch room.GameId {
	case 0: // Yahtzee //
		go yahtzee.Run(h.gamedb, h.hub, room)
	case 1: // LuckyNumber //
	}
}

func New(gamedb memorydb.GameDatabase, hub *model.WsHub) GameHandler {
	return &Games{gamedb, hub}
}
