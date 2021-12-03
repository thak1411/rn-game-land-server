package yahtzee

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"sort"
	"time"

	"github.com/thak1411/rn-game-land-server/memorydb"
	"github.com/thak1411/rn-game-land-server/model"
)

const (
	diceLen   = 5
	gameRound = 13

	// game response code //
	TypeSendRoomData    = 1
	TypeSendFieldDice   = 2
	TypeSendScore       = 3
	TypeSendTurn        = 4
	TypeSendRound       = 5
	TypeSendRerollCount = 6
	TypeGameEnd         = 8

	MsgTypeRollDice = 100
	MsgTypeGetScore = 101
)

var randSeed = rand.New(rand.NewSource(time.Now().UnixNano()))

type YahtzeeHandler struct {
	Turn        int             `json:"turn"`
	Round       int             `json:"round"`
	FieldDice   []int           `json:"fieldDice"`
	RerollCount int             `json:"rerollCount"`
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

func SendMessage(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId, msgType int, msg interface{}) {
	response := &YahtzeeResponse{}
	response.Code = 1000
	response.Message.Type = msgType
	response.Message.Data = msg

	narrowMsg, err := json.Marshal(response)
	if err != nil {
		log.Printf("error : %v", err)
		return
	}

	players := make([]int, 0)
	room, err := gamedb.GetRoom(roomId)
	if err != nil {
		log.Printf("error : %v", err)
		return
	}
	for _, v := range room.Player {
		if v.IsOnline {
			players = append(players, v.Id)
		}
	}

	narrowHandler := &model.NarrowcastHandler{
		Response: narrowMsg,
		Targets:  players,
	}
	hub.Narrowcast <- narrowHandler
}

func SendRoomData(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId int, data interface{}) {
	SendMessage(gamedb, hub, roomId, TypeSendRoomData, data)
}

func SendFieldDice(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId int, data interface{}) {
	SendMessage(gamedb, hub, roomId, TypeSendFieldDice, data)
}

func SendTurn(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId int, data interface{}) {
	SendMessage(gamedb, hub, roomId, TypeSendTurn, data)
}

func SendRound(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId int, data interface{}) {
	SendMessage(gamedb, hub, roomId, TypeSendRound, data)
}

func SendRerollCount(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId int, data interface{}) {
	SendMessage(gamedb, hub, roomId, TypeSendRerollCount, data)
}

func SendEndGame(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId int, data interface{}) {
	SendMessage(gamedb, hub, roomId, TypeGameEnd, data)
}

type ScoreResponse struct {
	Turn     int `json:"turn"`
	Score    int `json:"score"`
	ScoreKey int `json:"scoreKey"`
}

func SendScore(gamedb memorydb.GameDatabase, hub *model.WsHub, roomId, turn, score, scoreKey int) {
	response := &ScoreResponse{}
	response.Turn = turn
	response.Score = score
	response.ScoreKey = scoreKey

	SendMessage(gamedb, hub, roomId, TypeSendScore, response)
}

type ValidFunc func(interface{}) bool

func GetMessage(hub *model.WsHub, roomId, timeout int, message interface{}, isValid ValidFunc) int {
	tc := time.After(time.Duration(timeout) * time.Second)
	for {
		select {
		case msg := <-hub.GameMessage[roomId]:
			if err := json.Unmarshal(msg, message); err != nil {
				// filtered error -> no need //
				log.Printf("error : %v", err)
				continue
			}
			if isValid(message) {
				return 0
			}
		case <-tc:
			return -1
		}
	}
}

type BehaviorMessage struct {
	Code    int `json:"code"`
	Message struct {
		Id       int   `json:"id"`
		Type     int   `json:"type"`
		ScoreKey int   `json:"scoreKey"`
		Selected []int `json:"selected"`
	} `json:"message"`
}

func GetBehaviorMessage(hub *model.WsHub, roomId, timeout, playerId int, h *YahtzeeHandler) (*BehaviorMessage, int) {
	rollMsg := &BehaviorMessage{}
	isValid := func(message interface{}) bool {
		msg, ok := message.(*BehaviorMessage)
		if !ok {
			return false
		}
		if playerId != msg.Message.Id {
			return false
		}
		if msg.Message.Type == MsgTypeRollDice { // ReRoll Dice //
			if h.RerollCount >= 2 || len(msg.Message.Selected) < 1 {
				return false
			}
			dp := make(map[int]bool)
			for _, v := range msg.Message.Selected {
				if dp[v] || v < 0 || v >= diceLen {
					return false
				}
				dp[v] = true
			}
			return true
		} else if msg.Message.Type == MsgTypeGetScore { // Get Score //
			pack1 := msg.Message.ScoreKey >= 1 && msg.Message.ScoreKey <= 6
			pack2 := msg.Message.ScoreKey >= 9 && msg.Message.ScoreKey <= 15
			if !pack1 && !pack2 {
				return false
			}
			if _, ok := h.PlayerScore[h.Turn].Value[msg.Message.ScoreKey]; ok {
				return false
			}
			return true
		}
		return false
	}
	if ok := GetMessage(hub, roomId, timeout, rollMsg, isValid); ok != 0 {
		return nil, ok
	}
	return rollMsg, 0
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

func QueryScore(dice []int) (map[int]int, int, int, []int) {
	score := make(map[int]int)
	vlist := make([]int, 0)
	mx := 0
	sm := 0
	for _, v := range dice {
		score[v]++
		sm += v
		if mx < score[v] {
			mx = score[v]
		}
	}
	for _, v := range score {
		vlist = append(vlist, v)
	}
	sort.Ints(vlist)
	return score, mx, sm, vlist
}

func SetScore(h *YahtzeeHandler, scoreKey int) bool {
	var sc1, sc2, sc3 int
	scoreTable, mx, sm, vlist := QueryScore(h.FieldDice)

	// TODO: yahtzee 100 bonus score //
	if scoreKey >= 1 && scoreKey <= 6 {
		h.PlayerScore[h.Turn].Value[scoreKey] = scoreTable[scoreKey] * scoreKey
		h.PlayerScore[h.Turn].Value[7] += scoreTable[scoreKey] * scoreKey
		h.PlayerScore[h.Turn].Value[0] += scoreTable[scoreKey] * scoreKey
		if h.PlayerScore[h.Turn].Value[7] >= 0 {
			h.PlayerScore[h.Turn].Value[0] += 35
		}
	} else {
		switch scoreKey {
		case 9:
			if mx >= 3 {
				h.PlayerScore[h.Turn].Value[scoreKey] = sm
				h.PlayerScore[h.Turn].Value[0] += sm
			} else {
				h.PlayerScore[h.Turn].Value[scoreKey] = 0
			}
		case 10:
			if mx >= 4 {
				h.PlayerScore[h.Turn].Value[scoreKey] = sm
				h.PlayerScore[h.Turn].Value[0] += sm
			} else {
				h.PlayerScore[h.Turn].Value[scoreKey] = 0
			}
		case 11:
			if len(vlist) == 2 && vlist[0] == 2 && vlist[1] == 3 {
				h.PlayerScore[h.Turn].Value[scoreKey] = 25
				h.PlayerScore[h.Turn].Value[0] += 25
			} else {
				h.PlayerScore[h.Turn].Value[scoreKey] = 0
			}
		case 12:
			sc1 = 1
			sc2 = 1
			sc3 = 1
			for i := 1; i <= 4; i++ {
				sc1 *= scoreTable[i]
				sc2 *= scoreTable[i+1]
				sc3 *= scoreTable[i+2]
			}
			if sc1 > 0 || sc2 > 0 || sc3 > 0 {
				h.PlayerScore[h.Turn].Value[scoreKey] = 30
				h.PlayerScore[h.Turn].Value[0] += 30
			} else {
				h.PlayerScore[h.Turn].Value[scoreKey] = 0
			}
		case 13:
			sc1 = 1
			sc2 = 1
			for i := 1; i <= 5; i++ {
				sc1 *= scoreTable[i]
				sc2 *= scoreTable[i+1]
			}
			if sc1 > 0 || sc2 > 0 {
				h.PlayerScore[h.Turn].Value[scoreKey] = 40
				h.PlayerScore[h.Turn].Value[0] += 40
			} else {
				h.PlayerScore[h.Turn].Value[scoreKey] = 0
			}
		case 14:
			if len(vlist) == 1 {
				h.PlayerScore[h.Turn].Value[scoreKey] = 50
				h.PlayerScore[h.Turn].Value[0] += 50
			} else {
				h.PlayerScore[h.Turn].Value[scoreKey] = 0
			}
		case 15:
			h.PlayerScore[h.Turn].Value[scoreKey] = sm
			h.PlayerScore[h.Turn].Value[0] += sm
		default:
			return false
		}
	}
	return true
}

func IsYahtzee(dice []int) bool {
	for _, v := range dice {
		if v != dice[0] {
			return false
		}
	}
	return true
}

type YahtzeePlayerResult struct {
	Rank  int    `json:"rank"`
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type YahtzeeGameResult struct {
	ResultTable []*YahtzeePlayerResult `json:"resultTable"`
}

func Run(gamedb memorydb.GameDatabase, hub *model.WsHub, room *model.Room) {
	var _ = room.Option // game option
	var _ = room.Player // game player

	roomId := room.Id
	playerNum := len(room.Player)
	players := make([]int, playerNum)
	playerName := make([]string, playerNum)
	for i, v := range room.Player {
		players[i] = v.Id
		playerName[i] = v.Name
	}

	h := &YahtzeeHandler{
		Turn:        0,
		Round:       0,
		FieldDice:   make([]int, 0),
		RerollCount: 0,
		PlayerScore: make([]*YahtzeeScore, 0),
	}
	for i := 0; i < playerNum; i++ {
		h.PlayerScore = append(h.PlayerScore, &YahtzeeScore{
			Value: make(map[int]int),
		})
		h.PlayerScore[i].Value[0] = 0
		h.PlayerScore[i].Value[7] = -63
	}

	room.Data = h
	gamedb.SetRoomData(roomId, h)
	SendRoomData(gamedb, hub, roomId, h)

	turnTerm := 1000000

	for h.Round = 0; h.Round < gameRound; h.Round++ {
		SendRound(gamedb, hub, roomId, h.Round)
		for h.Turn = 0; h.Turn < playerNum; h.Turn++ {
			h.RerollCount = 0
			h.FieldDice = RollAllDice()

			if v, ok := h.PlayerScore[h.Turn].Value[15]; ok && v > 0 && IsYahtzee(h.FieldDice) {
				h.PlayerScore[h.Turn].Value[0] += 100
			}

			// TODO: Update: (round, turn) db handler //
			gamedb.SetRoomData(roomId, h) // TODO: Update to field dice handler //
			SendRerollCount(gamedb, hub, roomId, h.RerollCount)
			SendFieldDice(gamedb, hub, roomId, h.FieldDice)
			SendTurn(gamedb, hub, roomId, h.Turn)

			for {
				behaviorMsg, status := GetBehaviorMessage(hub, roomId, turnTerm, players[h.Turn], h)
				if status == 0 {
					if behaviorMsg.Message.Type == MsgTypeRollDice { // ReRoll Dice //
						if err := RollSelectedDice(h.FieldDice, behaviorMsg.Message.Selected); err != nil {
							// filtered error -> no need //
							log.Printf("error : %v", err)
							return
						}
						if v, ok := h.PlayerScore[h.Turn].Value[14]; ok && v > 0 && IsYahtzee(h.FieldDice) {
							h.PlayerScore[h.Turn].Value[0] += 100
						}
						h.RerollCount++
						gamedb.SetRoomData(roomId, h) // TODO: Update to field dice handler //
						SendFieldDice(gamedb, hub, roomId, h.FieldDice)
						SendRerollCount(gamedb, hub, roomId, h.RerollCount)
					} else if behaviorMsg.Message.Type == MsgTypeGetScore { // Get Score //
						scoreKey := behaviorMsg.Message.ScoreKey
						SetScore(h, scoreKey)
						gamedb.SetRoomData(roomId, h) // TODO: Update to field dice handler //
						SendScore(gamedb, hub, roomId, h.Turn, h.PlayerScore[h.Turn].Value[scoreKey], scoreKey)
						SendScore(gamedb, hub, roomId, h.Turn, h.PlayerScore[h.Turn].Value[0], 0)
						SendScore(gamedb, hub, roomId, h.Turn, h.PlayerScore[h.Turn].Value[7], 7)
						break
					}
				} else {
					// timeout //
					var _ = 1
					break
				}
			}
		}
	}
	result := &YahtzeeGameResult{
		ResultTable: nil,
	}
	for i, v := range h.PlayerScore {
		name := playerName[i]
		result.ResultTable = append(result.ResultTable, &YahtzeePlayerResult{
			Rank:  0,
			Name:  name,
			Score: v.Value[0],
		})
	}
	sort.Slice(result.ResultTable, func(i, j int) bool {
		return result.ResultTable[i].Score < result.ResultTable[j].Score
	})
	prevScore := -2
	rankId := 0
	for i, v := range result.ResultTable {
		if prevScore != v.Score {
			rankId++
		}
		result.ResultTable[i].Rank = rankId
		prevScore = v.Score
	}
	gamedb.SetGameEnd(roomId)
	SendEndGame(gamedb, hub, roomId, result)
}
