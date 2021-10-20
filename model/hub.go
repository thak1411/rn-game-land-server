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
	Join       chan *JoinForm
	Leave      chan *LeaveForm
}

type InviteForm struct {
	From       int    `json:"from"`
	RoomId     int    `json:"roomId"`
	TargetId   int    `json:"targetId"`
	TargetsId  []int  `json:"targetsId"`
	TargetName string `json:"targetName"`
}

type JoinForm struct {
	UserId    int   `json:"userId"`
	TargetsId []int `json:"targetId"`
}

type LeaveForm struct {
	UserId    int   `json:"userId"`
	TargetsId []int `json:"targetId"`
}
