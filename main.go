package main

import (
	"fmt"
	messages "github.com/tonnytg/messagebroker/entity/message"
	"github.com/tonnytg/messagebroker/pkg/servers/rabbit"
)

func main() {

	m := messages.NewMessage()
	m.GetMessage("Hello Friend!")
	fmt.Println(m.Data)

	// start connection with rabbit server
	conn, _ := rabbit.GetConn()

	for i:=0; i<10; i++ {
        // publish message to rabbit server
		m.GetMessage(fmt.Sprintf("Hello Friend %d", i))
		go rabbit.Publisher(&conn, m.Data)
    }
	//go rabbit.Publisher(&conn, m.Data)

	//rabbit.InitRabbitMQ()

	err := conn.StartConsumer("test-queue", "test-key", rabbit.Handler, 2)

	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	<-forever
}
