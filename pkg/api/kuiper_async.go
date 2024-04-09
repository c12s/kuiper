package api

import (
	"fmt"
	"log"

	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/c12s/magnetar/pkg/messaging/nats"
	natsgo "github.com/nats-io/nats.go"
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

func (c *KuiperAsyncClient) ReceiveConfig(standaloneHandler PutStandaloneConfigHandler, groupHandler PutConfigGroupHandler) error {
	err := c.subscriber.Subscribe(func(msg []byte, replySubject string) {
		// todo: prosiri tako da se vraca odgovor koji ce se vratiti nazad kuiper-u
		// o tome da li je task placed ili failed
		standaloneCmd := &ApplyStandaloneConfigCommand{}
		standaloneErr := standaloneCmd.Unmarshal(msg)
		if standaloneErr == nil {
			standaloneHandler(standaloneCmd)
			return
		}
		groupCmd := &ApplyConfigGroupCommand{}
		groupErr := groupCmd.Unmarshal(msg)
		if groupErr == nil {
			groupHandler(groupCmd)
			return
		}
		log.Println(standaloneErr)
		log.Println(groupErr)
	})
	return err
}

func (c *KuiperAsyncClient) GracefulStop() {
	err := c.subscriber.Unsubscribe()
	if err != nil {
		log.Println(err)
	}
}

type PutStandaloneConfigHandler func(cmd *ApplyStandaloneConfigCommand)
type PutConfigGroupHandler func(cmd *ApplyConfigGroupCommand)

func Subject(nodeId string) string {
	return fmt.Sprintf("%s.configs", nodeId)
}
