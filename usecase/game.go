package usecase

import (
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/model"
)

type GameUsecase interface {
	GetGameList() ([]model.Game, error)
	CreateRoom(int, int, string, string, string) (*model.Room, error)
	GetRoom(int) (*model.Room, error)
}

type GameUC struct {
	db database.GameDatabase
}

func (uc *GameUC) GetGameList() ([]model.Game, error) {
	return uc.db.GetGameList()
}

func (uc *GameUC) CreateRoom(owner, gameId int, name, option, ownerName string) (*model.Room, error) {
	return uc.db.CreateRoom(owner, gameId, name, option, ownerName)
}

func (uc *GameUC) GetRoom(roomId int) (*model.Room, error) {
	return uc.db.GetRoom(roomId)
}

func NewGame(db database.GameDatabase) GameUsecase {
	return &GameUC{db}
}
