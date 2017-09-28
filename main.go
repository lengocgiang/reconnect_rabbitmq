package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type RabbimtMQ struct {
	conn  *amqp.Connection
	chann *amqp.Channel
	done  chan error
}

var rabbit *RabbimtMQ

func NewRabbitMQ() *RabbimtMQ {
	return &RabbimtMQ{
		done: make(chan error),
	}
}

func initRabbitMQ() {
	rabbit = NewRabbitMQ()

	var err error

	rabbit.conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}

	rabbit.chann, err = rabbit.conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Printf("closing: %s", <-rabbit.conn.NotifyClose(make(chan *amqp.Error)))
		rabbit.done <- errors.New("Channel Closed")
	}()

	fmt.Println("Connect rabbitmq completed")

}

func main() {
	initRabbitMQ()
	startChan := make(chan int)

	go sendingRequest(rabbit, startChan)

	for i := 0; i < 100000; i++ {
		time.Sleep(1 * time.Second)
		startChan <- i
	}

}

func sendingRequest(rabbit *RabbimtMQ, startChan chan int) {
	fmt.Println("Sending request")
	pause := false
	var err error
	for {
		if !pause {
			select {
			case i := <-startChan:
				fmt.Println("Start :", i)
			case notifi := <-rabbit.done:
				fmt.Printf("Begin re-connect rabbitmq: %v \n", notifi)
				for {
					// auto reconnect in 3 seconds
					time.Sleep(3 * time.Second)
					rabbit.conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
					if err != nil {
						fmt.Printf("Re-connect rabbitmq failed: %v \n", err)
						continue
					}
					fmt.Printf("Re-connect rabbitmq succesfull \n", nil)
					pause = false
					break
				}
			}
		}
	}
}
