package broker

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type Conn struct {
	Channel *amqp.Channel
}

type ConfAMQP struct {
	Host         string
	Port         int
	User         string
	Password     string
	Vhost        string
	Exchange     string
	ExchangeType string
	RoutingKey   string
	Queue        string
	Durable      bool
	AutoDelete   bool
	Exclusive    bool
	NoWait       bool
	Args         amqp.Table
}

func NewConfigAMQO() *ConfAMQP {

	c := ConfAMQP{
		Host:         "localhost",
		Port:         1883,
		User:         "guest",
		Password:     "guest",
		Vhost:        "/",
		Exchange:     "events",
		ExchangeType: "direct",
		RoutingKey:   "test",
		Queue:        "test-queue",
		Durable:      true,
		AutoDelete:   false,
		Exclusive:    false,
		NoWait:       false,
		Args:         nil,
	}

	return &c
}

func Handler(d amqp.Delivery) bool {
	if d.Body == nil {
		fmt.Println("Error, no message body!")
		return false
	}
	fmt.Println(string(d.Body))
	return true
}

func Connect(c *ConfAMQP) (Conn, error) {

	serverURL := fmt.Sprintf("amqps://%s:%s@%s/%s", c.User, c.Password, c.Host, c.Vhost)

	conn, err := amqp.Dial(serverURL)
	if err != nil {
		log.Println("Error connecting to AMQP server: ", err)
		return Conn{}, err
	}
	// Success connected
	ch, err := conn.Channel()
	return Conn{
		Channel: ch,
	}, err
}

func Publish(conn *Conn, c *ConfAMQP, data []byte) error {

	conn.Channel.Publish(
		"events",
		c.RoutingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         data,
			DeliveryMode: amqp.Persistent,
		})
	return nil
}

func Subscriber(conn *Conn, c *ConfAMQP, handler func(d amqp.Delivery) bool, concurrency int) error {

	log.Println("Subscribing to queue: ", c.Queue)
	_, err := conn.Channel.QueueDeclare(c.Queue, c.Durable, c.AutoDelete, c.Exclusive, c.NoWait, c.Args)
	if err != nil {
		return err
	}

	// bind the queue to the routing key
	err = conn.Channel.QueueBind(c.Queue, c.RoutingKey, c.Exchange, c.NoWait, c.Args)
	if err != nil {
		log.Println("Error binding queue: ", err)
		return err
	}

	// prefetch 4x as many messages as we can handle at once
	prefetchCount := concurrency * 4
	err = conn.Channel.Qos(prefetchCount, 0, false)
	if err != nil {
		log.Println("Error setting QoS: ", err)
		return err
	}

	msgs, err := conn.Channel.Consume(
		c.Queue,     		// queue
		"",        // consumer
		false,      // auto-ack
		c.Exclusive, 		// exclusive
		false,      // no-local
		c.NoWait,    		// no-wait
		c.Args,      		// args
	)
	if err != nil {
		log.Println("Error connect to channel: ", err)
		return err
	}

	// create a goroutine for the number of concurrent threads requested
	for i := 0; i < concurrency; i++ {
		fmt.Printf("Processing messages on thread %v...\n", i)
		go func() {
			for msg := range msgs {
				// if tha handler returns true then ACK, else NACK
				// the message back into the rabbit queue for
				// another round of processing
				if handler(msg) {
					msg.Ack(false)
				} else {
					msg.Nack(false, true)
				}
			}
			fmt.Println("Rabbit consumer closed - critical Error")
			os.Exit(1)
		}()
	}
	return nil
}
