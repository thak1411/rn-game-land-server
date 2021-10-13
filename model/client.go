package model

import "github.com/gorilla/websocket"

/**
 * Client Object
 */
type ChatClient struct {
	Hub  *ChatHub
	Conn *websocket.Conn
	Send chan []byte
}
