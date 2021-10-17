package database

import "github.com/thak1411/rn-game-land-server/model"

type GameDatabase interface {
	GetGameList() ([]model.Game, error)
	CreateRoom(int, int, string, string, string) (*model.Room, error)
	GetRoom(int) (*model.Room, error)
}

type GameDB struct {
	GameList   []model.Game
	RoomList   map[int]*model.Room
	NextRoomId int
}

func (db *GameDB) GetGameList() ([]model.Game, error) {
	return db.GameList, nil
}

func (db *GameDB) CreateRoom(owner, gameId int, name, option, ownerName string) (*model.Room, error) {
	room := &model.Room{
		Id:     db.NextRoomId,
		Name:   name,
		Owner:  owner,
		GameId: gameId,
		Option: option,
		Player: []*model.Player{
			{
				Id:       owner,
				Name:     ownerName,
				IsOnline: false,
			},
		},
	}
	db.RoomList[db.NextRoomId] = room
	db.NextRoomId++
	return room, nil
}

func (db *GameDB) GetRoom(roomId int) (*model.Room, error) {
	return db.RoomList[roomId], nil
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
		RoomList:   make(map[int]*model.Room),
		NextRoomId: 1,
	}
}
