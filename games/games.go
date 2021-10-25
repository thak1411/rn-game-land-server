package games

import (
	"github.com/thak1411/rn-game-land-server/games/yahtzee"
	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/model"
)

type GameHandler interface {
	Run(*model.WsHub, *model.Room)
}

type Games struct {
	gamedb memorydb.GameDatabase
}

func (h *Games) Run(hub *model.WsHub, room *model.Room) {
	switch room.GameId {
	case 0: // Yahtzee //
		go yahtzee.Run(hub, room)
	case 1: // LuckyNumber //
	}
}

func New(gamedb memorydb.GameDatabase) GameHandler {
	return &Games{gamedb}
}
