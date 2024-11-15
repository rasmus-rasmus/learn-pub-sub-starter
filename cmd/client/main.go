package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/cmd/client/client"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	formatString := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(formatString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ. Error: %v", err)
	}
	numTries := 10
	var userName string
	for numTries > 0 {
		userName, err = gamelogic.ClientWelcome()
		if err != nil {
			fmt.Printf("%v\n", err)
			numTries--
		} else {
			break
		}
	}

	ch, err := conn.Channel()

	if err != nil {
		log.Fatalf("Could not create channel. Error: %v", err)
	}

	gameState := gamelogic.NewGameState(userName)

	// Subsribe to pause
	err = pubsub.SubscribeJSON(conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%v.%v", routing.PauseKey, userName),
		routing.PauseKey,
		1,
		client.HandlerPause(gameState))

	if err != nil {
		log.Fatalf("%v failed to subscribe to pause queue. Error: %v", userName, err)
	}

	// Subscribe to move
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%v.%v", routing.ArmyMovesPrefix, userName),
		fmt.Sprintf("%v.*", routing.ArmyMovesPrefix),
		1,
		client.HandlerMove(gameState, ch),
	)

	if err != nil {
		log.Fatalf("%v failed to subscribe to move queue. Error: %v", userName, err)
	}

	//Subscribe to war
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		fmt.Sprintf("%v.*", routing.WarRecognitionsPrefix),
		0,
		client.HandlerWar(gameState, ch),
	)

	if err != nil {
		log.Fatalf("%v failed to subscribe to move queue. Error: %v", userName, err)
	}

	breakOut := false

	for !breakOut {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "spawn":
			{
				err := gameState.CommandSpawn(input)
				if err != nil {
					fmt.Printf("Operation failed: %v\n", err)
				}
				break
			}
		case "move":
			{
				move, err := gameState.CommandMove(input)
				if err != nil {
					fmt.Printf("Operation failed: %v\n", err)
					break
				}
				err = pubsub.PublishJSON(
					ch,
					routing.ExchangePerilTopic,
					fmt.Sprintf("%v.%v", routing.ArmyMovesPrefix, userName),
					move,
				)
				if err != nil {
					log.Fatalf("Failed to publish move: %v. Error: %v", move, err)
				}
				log.Printf("Successfully published move: %v", move)
				break
			}
		case "status":
			{
				gameState.CommandStatus()
				break
			}
		case "help":
			{
				gamelogic.PrintClientHelp()
				break
			}
		case "spam":
			{
				if len(input) < 2 {
					fmt.Println("Please provide a spam number")
					break
				}
				numSpam, err := strconv.Atoi(input[1])
				if err != nil || numSpam <= 0 {
					fmt.Println("Please provide a positive integer")
					break
				}

				fmt.Printf("Spamming %v times\n", numSpam)

				for i := 0; i < numSpam; i++ {
					mallog := gamelogic.GetMaliciousLog()
					err = pubsub.PublishGob(
						ch,
						routing.ExchangePerilTopic,
						routing.GameLogSlug+"."+userName,
						routing.GameLog{
							CurrentTime: time.Now(),
							Message:     mallog,
							Username:    userName,
						},
					)
					if err != nil {
						fmt.Printf("Spamming failed. Error %v", err)
					}
				}
				break
			}
		case "quit":
			{
				gamelogic.PrintQuit()
				breakOut = true
				break
			}
		case "default":
			fmt.Println("Unknown operation.")
		}
	}
	log.Println("\nShutting down and closing connection to RabbitMQ.")
}
