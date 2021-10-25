package model

/**
 * Game Model Object
 *
 * TODO: Add Game Mode Option
 */
type Game struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	MinPlayer int    `json:"minPlayer"`
}

/**
 * Game Room Object
 */
type Room struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Start    bool      `json:"start"`
	Owner    int       `json:"owner"`
	GameId   int       `json:"gameId"`
	Option   string    `json:"option"`
	GameName string    `json:"gameName"`
	Player   []*Player `json:"player"`
}

type Player struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsOnline bool   `json:"isOnline"`
}
