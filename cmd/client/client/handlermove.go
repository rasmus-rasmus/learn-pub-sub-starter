package client

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) routing.AckType {
	return func(move gamelogic.ArmyMove) routing.AckType {
		defer fmt.Print("> ")
		moveOutCome := gs.HandleMove(move)
		switch moveOutCome {
		case gamelogic.MoveOutComeSafe:
			return routing.ACK
		case gamelogic.MoveOutcomeMakeWar:
			war := gamelogic.RecognitionOfWar{
				Attacker: move.Player,
				Defender: gs.GetPlayerSnap(),
			}
			err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%v.%v", routing.WarRecognitionsPrefix, gs.GetUsername()),
				war,
			)
			if err != nil {
				return routing.NACKREQUEUE
			}
			return routing.ACK
		case gamelogic.MoveOutcomeSamePlayer:
			return routing.NACKDISCARD
		default:
			return routing.NACKDISCARD
		}
	}
}
