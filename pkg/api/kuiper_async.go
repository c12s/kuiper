package api

import (
	context "context"
	"fmt"
	"log"

	"github.com/c12s/magnetar/pkg/messaging"
	"github.com/c12s/magnetar/pkg/messaging/nats"
	natsgo "github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

type KuiperAsyncClient struct {
	subscriber messaging.Subscriber
	publisher  messaging.Publisher
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
	publisher, err := nats.NewPublisher(conn)
	if err != nil {
		return nil, err
	}
	return &KuiperAsyncClient{
		subscriber: subscriber,
		publisher:  publisher,
	}, nil
}

func (c *KuiperAsyncClient) ReceiveConfig(standaloneHandler PutStandaloneConfigHandler, groupHandler PutConfigGroupHandler) error {
	return c.subscriber.Subscribe(func(ctx context.Context, msg []byte, replySubject string) {
		tracer := otel.Tracer("kuiper.async")
		ctx, span := tracer.Start(ctx, "ReceiveConfig")
		defer span.End()

		cmd := &ApplyConfigCommand{}
		err := proto.Unmarshal(msg, cmd)
		if err != nil {
			span.RecordError(err)
			log.Println(err)
			return
		}

		switch cmd.Type {
		case "standalone":
			config := &StandaloneConfig{}
			err := proto.Unmarshal(cmd.Config, config)
			if err != nil {
				span.RecordError(err)
				log.Println(err)
				return
			}
			err = standaloneHandler(config, cmd.Namespace, cmd.Strategy)
			reply := &ApplyConfigReply{Cmd: cmd}
			if err != nil {
				span.RecordError(err)
				reply.Status = TaskStatus_Failed
			} else {
				reply.Status = TaskStatus_Placed
			}
			msg, err := proto.Marshal(reply)
			if err != nil {
				span.RecordError(err)
				log.Println(err)
				return
			}
			err = c.publisher.Publish(ctx, msg, replySubject)
			if err != nil {
				span.RecordError(err)
				log.Println(err)
			}

		case "group":
			config := &ConfigGroup{}
			err := proto.Unmarshal(cmd.Config, config)
			if err != nil {
				span.RecordError(err)
				log.Println(err)
				return
			}
			err = groupHandler(config, cmd.Namespace, cmd.Strategy)
			reply := &ApplyConfigReply{Cmd: cmd}
			if err != nil {
				span.RecordError(err)
				reply.Status = TaskStatus_Failed
			} else {
				reply.Status = TaskStatus_Placed
			}
			msg, err := proto.Marshal(reply)
			if err != nil {
				span.RecordError(err)
				log.Println(err)
				return
			}
			err = c.publisher.Publish(ctx, msg, replySubject)
			if err != nil {
				span.RecordError(err)
				log.Println(err)
			}

		default:
			log.Printf("unknown cmd type %s", cmd.Type)
		}
	})
}

func (c *KuiperAsyncClient) GracefulStop() {
	err := c.subscriber.Unsubscribe()
	if err != nil {
		log.Println(err)
	}
}

type PutStandaloneConfigHandler func(config *StandaloneConfig, namespace, strategy string) error
type PutConfigGroupHandler func(config *ConfigGroup, namespace, strategy string) error

func Subject(nodeId string) string {
	return fmt.Sprintf("%s.configs", nodeId)
}
