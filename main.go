package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"messagebroker/pkg/rabbit"
	"os"
	"time"
)

func main() {
	// I used https://api.cloudamqp.com in free tier for this example
	user := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")
	host := os.Getenv("RABBITMQ_HOST")
	port := os.Getenv("RABBITMQ_PORT")
	vhost := os.Getenv("RABBITMQ_VHOST")
	database := os.Getenv("RABBITMQ_DATABASE")

	conn, err := rabbit.GetConn(fmt.Sprintf("amqp://%s:%s@%s/%s", user, password, host, port, vhost, database))
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
