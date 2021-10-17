package database

import "github.com/thak1411/rn-game-land-server/model"

type GameDatabase interface {
	GetGameList() ([]model.Game, error)
}

type GameDB struct {
	GameList []model.Game
}

func (db *GameDB) GetGameList() ([]model.Game, error) {
	return db.GameList, nil
}

func NewGame() GameDatabase {
	return &GameDB{
		GameList: []model.Game{
			{
				Id:        0,
				Name:      "Yahtzee",
				MinPlayer: 2,
			},
		},
	}
}
