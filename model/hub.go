package model

/**
 * WebSocket Clients Hub
 * Managing Client & Sending Message
 */
type Hub struct {
	Clients    map[*Client]string
	Register   chan *Client
	UnRegister chan *Client
}
