package client

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func HandlerPause(gs *gamelogic.GameState) func(routing.PlayingState) routing.AckType {
	return func(state routing.PlayingState) routing.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(state)
		return routing.ACK
	}
}
