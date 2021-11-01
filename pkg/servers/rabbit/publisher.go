package rabbit

import (
	"fmt"
	"time"
)

func Publisher(conn *Conn, msg interface{}) {
	for {
		time.Sleep(time.Second)
		m := fmt.Sprintf("{\"message\":\"%s\"}", msg)
		conn.Publish("test-key", []byte(m))
	}
}
