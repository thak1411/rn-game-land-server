package model

import "github.com/gorilla/websocket"

type WsDefaultMessage struct {
	Code int `json:"code"`
	// Message interface{} `json:"message"`
}

type WsClient struct {
	Hub  *WsHub
	Conn *websocket.Conn
	Send chan []byte
	WsUser
}

type WsHub struct {
	Clients      map[int]*WsClient
	Register     chan *WsClient
	UnRegister   chan *WsClient
	Broadcast    chan []byte
	Narrowcast   chan *NarrowcastHandler
	BroadcastLog [][]byte
}

type WsUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	RoomId   int    `json:"roomId"`
	Username string `json:"username"`
}

type NarrowcastHandler struct {
	Targets  []int
	Response []byte
}

const (
	// config //
	BC_LOG_LEN = 20

	// status code //
	RES_BROADCAST = 100

	// message code //
	REQ_BROADCAST = 90
	REQ_INVITE    = 50
	REQ_JOIN      = 51
	REQ_REJECT    = 52
	REQ_START     = 20
)
