package model

/**
 * WebSocket Chatting Clients Hub
 * Managing Client & Sending Message
 */
type ChatHub struct {
	Clients    map[*ChatClient]bool
	Register   chan *ChatClient
	UnRegister chan *ChatClient
	Broadcast  chan []byte
	LastLog    []string
}

/**
 * WebSocket Notice Clients Hub
 */
type NoticeHub struct {
	Clients    map[int]*NoticeClient
	Register   chan *NoticeClient
	UnRegister chan *NoticeClient
	Invite     chan *InviteForm
	InviteLog  map[int][]*InviteForm
}

type InviteForm struct {
	From     int `json:"from"`
	RoodId   int `json:"roodId"`
	TargetId int `json:"targetId"`
}
