package model

/**
 * WebSocket Clients Hub
 * Managing Client & Sending Message
 */
type ChatHub struct {
	Clients    map[*ChatClient]ChatUser
	Register   chan *ChatClient
	UnRegister chan *ChatClient
	Broadcast  chan []byte
	LastLog    []string
}
