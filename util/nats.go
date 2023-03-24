package util

import "github.com/nats-io/nats.go"

func Conn(natsAddr string) (*nats.Conn, error) {
	return nats.Connect(natsAddr)
}
