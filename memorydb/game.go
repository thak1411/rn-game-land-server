package memorydb

import (
	"errors"
	"math/rand"

	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rnjson"
)

type GameDatabase interface {
	GetGameList() ([]model.Game, error)
	GetGameName(int) (string, error)
	CreateRoom(int, int, string, string, string, string) (*model.Room, error)
	GetRoom(int) (*model.Room, error)
	SetUserOnline(int, int) (bool, error)
	SetUserOffline(int) (bool, *model.Room, error)
	AppendRoomPlayer(int, int, string) (bool, error)
	DeleteRoomPlayer(int, int) (bool, error)
	GetInviteMessage(int) ([][]byte, error)
	AppendInviteMessage(int, []byte) (bool, error)
	DeleteInviteMessage(int, int) (int, error)
	SetGameStart(int) (bool, error)
	SetGameEnd(int) (bool, error)
	ShufflePlayer(int) error
	SetRoomData(int, interface{}) error
	GetRoomPlayer(int) ([]*model.Player, error)
}

type GameDB struct {
	GameList      []model.Game
	RoomList      map[int]*model.Room
	NextRoomId    int
	InviteMessage map[int][][]byte
}

func (db *GameDB) GetGameList() ([]model.Game, error) {
	return db.GameList, nil
}

func (db *GameDB) GetGameName(gameId int) (string, error) {
	if gameId < 0 || gameId >= len(db.GameList) {
		return "", nil
	}
	return db.GameList[gameId].Name, nil
}

func (db *GameDB) CreateRoom(owner, gameId int, name, gameName, option, ownerName string) (*model.Room, error) {
	room := &model.Room{
		Id:       db.NextRoomId,
		Data:     nil,
		Name:     name,
		Owner:    owner,
		Start:    false,
		GameId:   gameId,
		Option:   option,
		GameName: gameName,
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
	room, ok := db.RoomList[roomId]
	if ok {
		return room, nil
	}
	return nil, errors.New("no room")
}

func (db *GameDB) SetUserOnline(roomId, userId int) (bool, error) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return false, err
	}
	for i, v := range room.Player {
		if v.Id == userId {
			room.Player[i].IsOnline = true
			db.RoomList[roomId] = room
			return true, nil
		}
	}
	return false, nil
}

func (db *GameDB) SetUserOffline(userId int) (bool, *model.Room, error) {
	for i, v := range db.RoomList {
		for j, w := range v.Player {
			if w.Id == userId && w.IsOnline {
				db.RoomList[i].Player[j].IsOnline = false
				return true, v, nil
			}
		}
	}
	return false, nil, nil
}

func (db *GameDB) AppendRoomPlayer(roomId, userId int, userName string) (bool, error) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return false, err
	}
	for _, v := range room.Player {
		if v.Id == userId {
			return false, nil
		}
	}
	room.Player = append(room.Player, &model.Player{
		Id:       userId,
		Name:     userName,
		IsOnline: false,
	})
	db.RoomList[roomId] = room
	return false, nil
}

func (db *GameDB) DeleteRoomPlayer(roomId, userId int) (bool, error) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return false, err
	}
	for i, v := range room.Player {
		if v.Id == userId {
			room.Player = append(room.Player[:i], room.Player[i+1:]...)
			db.RoomList[roomId] = room
			return true, nil
		}
	}
	return false, nil
}

func (db *GameDB) GetInviteMessage(userId int) ([][]byte, error) {
	msg, ok := db.InviteMessage[userId]
	if !ok {
		return nil, nil
	}
	return msg, nil
}

func (db *GameDB) AppendInviteMessage(userId int, message []byte) (bool, error) {
	db.InviteMessage[userId] = append(db.InviteMessage[userId], message)
	return true, nil
}

func (db *GameDB) DeleteInviteMessage(userId, roomId int) (int, error) {
	msg, ok := db.InviteMessage[userId]
	if !ok {
		return 0, nil
	}
	newMsg := [][]byte{}
	deleteCount := 0
	for _, v := range msg {
		message, err := rnjson.Unmarshal(string(v))
		if err != nil {
			return 0, err
		}
		msgRoomId, ok := rnjson.Get(message, "message.roomId").(float64)
		if !ok {
			return 0, errors.New("invalid message")
		}
		if int(msgRoomId) == roomId {
			deleteCount++
			continue
		}
		newMsg = append(newMsg, v)
	}
	db.InviteMessage[userId] = newMsg
	return deleteCount, nil
}

func (db *GameDB) SetGameStart(roomId int) (bool, error) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return false, err
	}
	res := !room.Start
	room.Start = true
	db.RoomList[roomId] = room
	return res, nil
}

func (db *GameDB) SetGameEnd(roomId int) (bool, error) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return false, err
	}
	res := !room.Start
	room.Start = false
	db.RoomList[roomId] = room
	return res, nil
}

func (db *GameDB) ShufflePlayer(roomId int) error {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return err
	}
	rand.Shuffle(len(room.Player), func(i, j int) { room.Player[i], room.Player[j] = room.Player[j], room.Player[i] })
	db.RoomList[roomId] = room
	return nil
}

func (db *GameDB) SetRoomData(roomId int, data interface{}) error {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return err
	}
	room.Data = data
	db.RoomList[roomId] = room
	return nil
}

func (db *GameDB) GetRoomPlayer(roomId int) ([]*model.Player, error) {
	room, err := db.GetRoom(roomId)
	if err != nil {
		return nil, err
	}
	return room.Player, nil
}

var gameDB GameDatabase = nil

func NewGame() GameDatabase {
	if gameDB == nil {
		gameDB = &GameDB{
			GameList: []model.Game{
				{
					Id:        0, // Auto Increase //
					Name:      "Yahtzee",
					MinPlayer: 2,
				},
			},
			RoomList:      make(map[int]*model.Room),
			NextRoomId:    1,
			InviteMessage: make(map[int][][]byte),
		}
	}
	return gameDB
}
