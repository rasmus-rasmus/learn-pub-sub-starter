package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	formatString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(formatString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ. Error: %v", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not create channel. Error: %v", err)
	}
	defer conn.Close()
	fmt.Println("Server successfully connected to RabbitMQ")

	err = pubsub.SubscribeGob(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		fmt.Sprintf("%v.*", routing.GameLogSlug),
		0,
		func(gl routing.GameLog) routing.AckType {
			defer fmt.Print("> ")
			gamelogic.WriteLog(gl)
			return routing.ACK
		},
	)

	if err != nil {
		log.Fatalf("Could not create queue. Error: %v", err)
	}

	gamelogic.PrintServerHelp()

	breakOut := false

	for !breakOut {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			{
				err = pubsub.PublishJSON(ch,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{IsPaused: true})
				if err != nil {
					log.Printf("Could not publish pause message. Error: %v", err)
				}
				log.Println("Pause message sent.")
				break
			}
		case "resume":
			{
				err = pubsub.PublishJSON(ch,
					routing.ExchangePerilDirect,
					routing.PauseKey,
					routing.PlayingState{IsPaused: false})
				if err != nil {
					log.Printf("Could not publish pause message. Error: %v", err)
				}
				log.Println("Resume message sent.")
				break
			}
		case "quit":
			{
				fmt.Println("Exiting...")
				breakOut = true
				break
			}
		default:
			fmt.Println("Unknown command")
		}
	}
	log.Println("Shutting down and closing connection to RabbitMQ.")
}
