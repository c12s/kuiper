package startup

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func NewNatsConn(address string) (*nats.Conn, error) {
	return nats.Connect(fmt.Sprintf("nats://%s", address))
}
