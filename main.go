package main

import (
	"flag"
	"fmt"
	"messagebroker/pkg/rabbit"
	"time"
)

func main() {

	fmt.Println("Reading message...")
	message := flag.String("message", "", "Message to send")
	flag.Parse()
	fmt.Println(*message)
	time.Sleep(time.Second * 1)

	rabbit.InitRabbitMQ()
}
