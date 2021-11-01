package broker_test

import (
	"github.com/tonnytg/messagebroker"
	"os"
	"testing"
	"time"
)

func TestConnectBroker(t *testing.T) {

	// create new configuration
	c := broker.NewConfigAMQO()

	// set new values to configuration
	c.User = os.Getenv("RABBITMQ_USER")
	c.Password = os.Getenv("RABBITMQ_PASSWORD")
	c.Host = os.Getenv("RABBITMQ_HOST")
	c.Vhost = os.Getenv("RABBITMQ_DATABASE")
	c.Queue = "test"
	c.RoutingKey = os.Getenv("RABBITMQ_KEY")

	_, err := broker.Connect(c)
	if err != nil {
		t.Log("Connected to broker")
    }

	conn, err := broker.Connect(c)
	broker.Publish(&conn, c, []byte("go test"))
	time.Sleep(time.Second * 1)
	err = broker.Subscriber(&conn, c, broker.Handler, 1)
	time.Sleep(time.Second * 1)
	if err != nil {
		t.Log("Error to subscriber")
		panic(err)
	}
}
