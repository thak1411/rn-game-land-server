package usecase

import (
	"github.com/thak1411/rn-game-land-server/database"
	"github.com/thak1411/rn-game-land-server/model"
)

type GameUsecase interface {
	GetGameList() ([]model.Game, error)
	CreateRoom(int, int, string, string, string) (*model.Room, error)
	GetRoom(int, int) (*model.Room, error)
}

type GameUC struct {
	db database.GameDatabase
}

func (uc *GameUC) GetGameList() ([]model.Game, error) {
	return uc.db.GetGameList()
}

func (uc *GameUC) CreateRoom(owner, gameId int, name, option, ownerName string) (*model.Room, error) {
	gameName, err := uc.db.GetGameName(gameId)
	if err != nil {
		return nil, err
	}
	if gameName == "" {
		return nil, nil
	}
	return uc.db.CreateRoom(owner, gameId, name, gameName, option, ownerName)
}

func (uc *GameUC) GetRoom(userId, roomId int) (*model.Room, error) {
	room, err := uc.db.GetRoom(roomId)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, nil
	}
	for _, p := range room.Player {
		if p.Id == userId {
			return room, nil
		}
	}
	return nil, nil
}

func NewGame(db database.GameDatabase) GameUsecase {
	return &GameUC{db}
}
