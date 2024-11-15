package client

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func HandlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) routing.AckType {
	return func(war gamelogic.RecognitionOfWar) routing.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(war)
		logMsg := ""
		if warOutcome == gamelogic.WarOutcomeOpponentWon || warOutcome == gamelogic.WarOutcomeYouWon {
			logMsg = fmt.Sprintf("%v won a war against %v", winner, loser)
		} else if warOutcome == gamelogic.WarOutcomeDraw {
			logMsg = fmt.Sprintf("A war between %v and %v resulted in a draw", winner, loser)
		}

		if len(logMsg) > 0 {

			err := pubsub.PublishGob(publishCh,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%v.%v", routing.GameLogSlug, gs.GetUsername()),
				routing.GameLog{
					CurrentTime: time.Now(),
					Message:     logMsg,
					Username:    gs.GetUsername(),
				})
			if err != nil {
				return routing.NACKREQUEUE
			}
		}

		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return routing.NACKREQUEUE
		case gamelogic.WarOutcomeNoUnits:
			return routing.NACKDISCARD
		case gamelogic.WarOutcomeOpponentWon:
			return routing.ACK
		case gamelogic.WarOutcomeYouWon:
			return routing.ACK
		case gamelogic.WarOutcomeDraw:
			return routing.ACK
		default:
			fmt.Printf("Unknown outcome type: %v\n", warOutcome)
			return routing.NACKDISCARD
		}
	}
}
