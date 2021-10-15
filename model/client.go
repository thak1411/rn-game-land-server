package model

import "github.com/gorilla/websocket"

/**
 * Client Object
 */
type ChatClient struct {
	Hub  *ChatHub
	Conn *websocket.Conn
	Send chan []byte
	ChatUser
}

/**
 * Client User Info Object
 */
type ChatUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}
