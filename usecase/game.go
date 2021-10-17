package usecase

import (
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/model"
)

type GameUsecase interface {
	GetGameList() ([]model.Game, error)
}

type GameUC struct {
	db database.GameDatabase
}

func (uc *GameUC) GetGameList() ([]model.Game, error) {
	return uc.db.GetGameList()
}

func NewGame(db database.GameDatabase) GameUsecase {
	return &GameUC{db}
}
