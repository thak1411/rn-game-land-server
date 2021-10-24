package database

// import (
// 	"errors"

// 	"github.com/thak1411/rn-game-land-server/model"
// )

// type GameDatabase interface {
// 	GetGameList() ([]model.Game, error)
// 	GetGameName(int) (string, error)
// 	CreateRoom(int, int, string, string, string, string) (*model.Room, error)
// 	GetRoom(int) (*model.Room, error)
// 	SetUserOnline(int, int) (bool, error)
// 	SetUserOffline(int) (bool, *model.Room, error)
// 	AppendRoomPlayer(int, int, string) (bool, error)
// }

// type GameDB struct {
// 	GameList   []model.Game
// 	RoomList   map[int]*model.Room
// 	NextRoomId int
// }

// func (db *GameDB) GetGameList() ([]model.Game, error) {
// 	return db.GameList, nil
// }

// func (db *GameDB) GetGameName(gameId int) (string, error) {
// 	if gameId < 0 || gameId >= len(db.GameList) {
// 		return "", nil
// 	}
// 	return db.GameList[gameId].Name, nil
// }

// func (db *GameDB) CreateRoom(owner, gameId int, name, gameName, option, ownerName string) (*model.Room, error) {
// 	room := &model.Room{
// 		Id:       db.NextRoomId,
// 		Name:     name,
// 		Owner:    owner,
// 		GameId:   gameId,
// 		Option:   option,
// 		GameName: gameName,
// 		Player: []*model.Player{
// 			{
// 				Id:       owner,
// 				Name:     ownerName,
// 				IsOnline: false,
// 			},
// 		},
// 	}
// 	db.RoomList[db.NextRoomId] = room
// 	db.NextRoomId++
// 	return room, nil
// }

// func (db *GameDB) GetRoom(roomId int) (*model.Room, error) {
// 	room, ok := db.RoomList[roomId]
// 	if ok {
// 		return room, nil
// 	}
// 	return nil, errors.New("no room")
// }

// func (db *GameDB) SetUserOnline(roomId, userId int) (bool, error) {
// 	room, err := db.GetRoom(roomId)
// 	if err != nil {
// 		return false, err
// 	}
// 	for i, v := range room.Player {
// 		if v.Id == userId {
// 			room.Player[i].IsOnline = true
// 			db.RoomList[roomId] = room
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// func (db *GameDB) SetUserOffline(userId int) (bool, *model.Room, error) {
// 	for i, v := range db.RoomList {
// 		for j, w := range v.Player {
// 			if w.Id == userId && w.IsOnline {
// 				db.RoomList[i].Player[j].IsOnline = false
// 				return true, v, nil
// 			}
// 		}
// 	}
// 	return false, nil, nil
// }

// func (db *GameDB) AppendRoomPlayer(roomId, userId int, userName string) (bool, error) {
// 	room, err := db.GetRoom(roomId)
// 	if err != nil {
// 		return false, err
// 	}
// 	for _, v := range room.Player {
// 		if v.Id == userId {
// 			return false, nil
// 		}
// 	}
// 	room.Player = append(room.Player, &model.Player{
// 		Id:       userId,
// 		Name:     userName,
// 		IsOnline: false,
// 	})
// 	db.RoomList[roomId] = room
// 	return false, nil
// }

// var gameDB GameDatabase = nil

// func NewGame() GameDatabase {
// 	if gameDB == nil {
// 		gameDB = &GameDB{
// 			GameList: []model.Game{
// 				{
// 					Id:        0, // Auto Increase //
// 					Name:      "Yahtzee",
// 					MinPlayer: 2,
// 				},
// 			},
// 			RoomList:   make(map[int]*model.Room),
// 			NextRoomId: 1,
// 		}
// 	}
// 	return gameDB
// }
