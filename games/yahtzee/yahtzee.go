package yahtzee

import (
	"fmt"

	"github.com/thak1411/rn-game-land-server/model"
)

func Run(hub *model.WsHub, room *model.Room) {
	var _ = room.Option // game option
	var _ = room.Player // game player

	fmt.Println("start Yahtzee")
}
