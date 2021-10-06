package usecase

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

type ChatUsecase interface {
	SocketTest(*websocket.Conn) error
}

type ChatUC struct{}

func (uc *ChatUC) SocketTest(conn *websocket.Conn) error {
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			err := errors.New("internal server error")
			return err
		}
		fmt.Println("CONN", conn)
		fmt.Println("RCV", msgType)
		fmt.Println("RCV", msg)
		fmt.Println()
		if err := conn.WriteMessage(msgType, msg); err != nil {
			err := errors.New("internal server error")
			return err
		}
	}
	// return nil
}

func NewChat() ChatUsecase {
	return &ChatUC{}
}
