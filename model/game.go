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
	Id     int
	Name   string
	Owner  int
	GameId int
	Option string
	Player []*Player
}

type Player struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	IsOnline bool   `json:"isOnline"`
}
