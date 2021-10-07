package model

import "github.com/gorilla/websocket"

/**
 * Client Object
 */
type Client struct {
	Hub  *Hub
	Conn *websocket.Conn
	Send chan []byte
}
