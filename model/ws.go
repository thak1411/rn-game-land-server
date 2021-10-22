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
	BroadcastLog [][]byte
}

type WsUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	RoomId   int    `json:"roomId"`
	Username string `json:"username"`
}

const (
	BroadcastLogLen = 20
)
