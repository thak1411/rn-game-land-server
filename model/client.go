package model

import "github.com/gorilla/websocket"

/**
 * Chatting Client Object
 */
type ChatClient struct {
	Hub  *ChatHub
	Conn *websocket.Conn
	Send chan []byte
	ChatUser
}

/**
 * Chatting Client User Info Object
 */
type ChatUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

/**
 * Notice Client Object
 */
type NoticeClient struct {
	Hub  *NoticeHub
	Conn *websocket.Conn
	Send chan []byte
	NoticeUser
}

/**
 * Notice Client User Info Object
 */
type NoticeUser struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	RoomId   int    `json:"roomId"`
	Username string `json:"username"`
}
