package api

import (
	"fmt"
	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/c12s/magnetar/pkg/messaging/nats"
	natsgo "github.com/nats-io/nats.go"
	"log"
)

type KuiperAsyncClient struct {
	subscriber messaging.Subscriber
}

func NewKuiperAsyncClient(address, nodeId string) (*KuiperAsyncClient, error) {
	conn, err := natsgo.Connect(fmt.Sprintf("nats://%s", address))
	if err != nil {
		return nil, err
	}
	subscriber, err := nats.NewSubscriber(conn, Subject(nodeId), nodeId)
	if err != nil {
		return nil, err
	}
	return &KuiperAsyncClient{
		subscriber: subscriber,
	}, nil
}

func (c *KuiperAsyncClient) ReceiveConfig(handler PutConfigHandler) error {
	err := c.subscriber.Subscribe(func(msg []byte, replySubject string) {
		cmd := &ApplyConfigCommand{}
		err := cmd.Unmarshal(msg)
		if err != nil {
			log.Println(err)
			return
		}
		handler(cmd)
	})
	return err
}

func (c *KuiperAsyncClient) GracefulStop() {
	err := c.subscriber.Unsubscribe()
	if err != nil {
		log.Println(err)
	}
}

type PutConfigHandler func(cmd *ApplyConfigCommand)

func Subject(nodeId string) string {
	return fmt.Sprintf("%s.configs", nodeId)
}
