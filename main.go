package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"messagebroker/rabbit"
	"time"
)

func main() {
	conn, err := rabbit.GetConn("amqp://xreavpmp:Gu-USkeppkYYeeWfMOQ4ynEoqkyXb6Y-@elk.rmq2.cloudamqp.com/xreavpmp")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(time.Second)
			conn.Publish("test-key", []byte(`{"message":"test1"}`))
		}
	}()

	err = conn.StartConsumer("test-queue", "test-key", handler, 2)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	<-forever
}

func handler(d amqp.Delivery) bool {
	if d.Body == nil {
		fmt.Println("Error, no message body!")
		return false
	}
	fmt.Println(string(d.Body))
	return true
}
