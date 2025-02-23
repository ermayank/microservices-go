package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"listener-service/event"
	"log"
	"math"
	"os"
	"time"
)

func main() {
	//Connect to Rabbitmq
	rabbitConn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer rabbitConn.Close()

	//Start Listen to messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// Create Consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// Watch the queue and consume events
	err = consumer.Listen([]string{
		"log.INFO", "log.WARNING", "log.ERROR",
	})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	//Don't continue until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ is not ready")
			counts++
		} else {
			log.Println("Connected to RabbitMQ")
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("RabbitMQ failed to connect")
			return nil, err
		}

		backOff = time.Duration(math.Pow(2, float64(backOff))) * time.Second
		log.Println("Backing Off !")
		time.Sleep(backOff)
		continue
	}
	return connection, nil
}
