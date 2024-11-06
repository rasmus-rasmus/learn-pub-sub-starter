package client

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func HandlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) routing.AckType {
	return func(move gamelogic.ArmyMove) routing.AckType {
		defer fmt.Print("> ")
		moveOutCome := gs.HandleMove(move)
		switch moveOutCome {
		case gamelogic.MoveOutComeSafe, gamelogic.MoveOutcomeMakeWar:
			return routing.ACK
		case gamelogic.MoveOutcomeSamePlayer:
			return routing.NACKDISCARD
		default:
			return routing.NACKDISCARD
		}
	}
}
