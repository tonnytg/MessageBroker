package main

import (
	"fmt"
	"github.com/tonnytg/messagebroker"
	"github.com/tonnytg/messagebroker/entity/message"
	"log"
	"os"
)

func main() {

	c := broker.NewConfigAMQO()

	c.User = os.Getenv("RABBITMQ_USER")
	c.Password = os.Getenv("RABBITMQ_PASSWORD")
	c.Host = os.Getenv("RABBITMQ_HOST")
	c.Vhost = os.Getenv("RABBITMQ_DATABASE")
	c.Queue = os.Getenv("RABBITMQ_QUEUE")
	c.RoutingKey = os.Getenv("RABBITMQ_KEY")
	c.Exchange = os.Getenv("RABBIT_EXCHANGE")

	m := messages.NewMessage()
	m.GetMessage("Hello Friend!")
	fmt.Println(m.Data)

	// start connection with rabbit server
	conn, err := broker.Connect(c)
	if err != nil {
		log.Println("ERROR: AMQP Server ", err)
	}

	for i := 0; i < 10; i++ {
		// publish message to rabbit server
		m.GetMessage(fmt.Sprintf("Hello Friend %d", i))
		go broker.Publish(&conn, c, []byte("test"))
	}

	err = broker.Subscriber(&conn, c, broker.Handler, 2)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	<-forever
}
