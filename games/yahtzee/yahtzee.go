package yahtzee

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/model"
)

const (
	diceLen   = 5
	gameRound = 13

	// game response code //
	TypeSendRoomData  = 1
	TypeSendFieldDice = 2
)

var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

type YahtzeeHandler struct {
	Turn        int             `json:"turn"`
	Round       int             `json:"round"`
	FieldDice   []int           `json:"fieldDice"`
	PlayerScore []*YahtzeeScore `json:"playerScore"`
}

type YahtzeeScore struct {
	Value map[int]int `json:"value"`
}

type YahtzeeResponse struct {
	Code    int `json:"code"`
	Message struct {
		Type int         `json:"type"`
		Data interface{} `json:"data"`
	} `json:"message"`
}

func SendMessage(hub *model.WsHub, players []int, msgType int, msg interface{}) {
	response := &YahtzeeResponse{}
	response.Code = 1000
	response.Message.Type = msgType
	response.Message.Data = msg

	narrowMsg, err := json.Marshal(response)
	if err != nil {
		log.Printf("error : %v", err)
		return
	}

	narrowHandler := &model.NarrowcastHandler{
		Response: narrowMsg,
		Targets:  players,
	}
	hub.Narrowcast <- narrowHandler
}

func SendRoomData(hub *model.WsHub, players []int, data interface{}) {
	SendMessage(hub, players, TypeSendRoomData, data)
}

func SendFieldDice(hub *model.WsHub, players []int, data interface{}) {
	SendMessage(hub, players, TypeSendFieldDice, data)
}

func GetMessage(hub *model.WsHub, roomId, timeout int) (bool, []byte) {
	select {
	case msg := <-hub.GameMessage[roomId]:
		return true, msg
	case <-time.After(time.Duration(timeout) * time.Second):
		return false, nil
	}
}

func RollDice() int {
	return randSeed.Intn(6) + 1
}

func RollAllDice() []int {
	dice := make([]int, diceLen)
	for i := 0; i < diceLen; i++ {
		dice[i] = RollDice()
	}
	return dice
}

func RollSelectedDice(dice, selected []int) error {
	dp := make(map[int]int)
	for _, v := range selected {
		if dp[v] == 1 {
			return errors.New("duplicated select in rolling dice")
		}
		if v < 0 || v >= diceLen {
			return errors.New("invalid index to rolling dice")
		}
		dice[v] = RollDice()
		dp[v] = 1
	}
	return nil
}

func Run(gamedb memorydb.GameDatabase, hub *model.WsHub, room *model.Room) {
	var _ = room.Option // game option
	var _ = room.Player // game player

	roomId := room.Id
	playerNum := len(room.Player)
	players := make([]int, playerNum)
	for i, v := range room.Player {
		players[i] = v.Id
	}

	h := &YahtzeeHandler{
		Turn:        0,
		Round:       0,
		FieldDice:   make([]int, 0),
		PlayerScore: make([]*YahtzeeScore, 0),
	}
	for i := 0; i < playerNum; i++ {
		h.PlayerScore = append(h.PlayerScore, &YahtzeeScore{
			Value: make(map[int]int),
		})
	}

	room.Data = h
	gamedb.SetRoomData(roomId, h)
	SendRoomData(hub, players, h)

	turnTerm := 10

	for h.Round = 0; h.Round < gameRound; h.Round++ {
		for h.Turn = 0; h.Turn < playerNum; h.Turn++ {
			h.FieldDice = RollAllDice()
			// TODO: Update: (round, turn) db handler //
			gamedb.SetRoomData(roomId, h) // TODO: Update to field dice handler //
			SendFieldDice(hub, players, h.FieldDice)

			res, msg := GetMessage(hub, roomId, turnTerm)
			if res {
				fmt.Println("msg", msg)
			} else {
				fmt.Println("timeout")
			}
		}
	}
}
