package main

import (
	"fmt"
	"log"

	"time"

	"github.com/streadway/amqp"
)

var amqpUri = "amqp://guest:guest@localhost:5672/"
var rabbitCloseError chan *amqp.Error

type RabbitMQ struct {
}

func main() {
	for {
		fmt.Println("Begin create new connection to rabbitmq in 5 second.")
		time.Sleep(5 * time.Second)

		conn := connect()
		//done := make(chan error)

		go func() {
			log.Printf("closing here %s", <-conn.NotifyClose(make(chan *amqp.Error)))
			//done <- errors.New("Channel Close")
			conn = connect()
		}()

		log.Println("got connection, creating channel")
	}
}

func connect() *amqp.Connection {
	for {
		conn, err := amqp.Dial(amqpUri)
		if err == nil {
			return conn
		}
		log.Println(err)
		log.Println("Trying to reconnect to rabbitmq in 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
