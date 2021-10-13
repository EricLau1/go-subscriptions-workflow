package rmq

import (
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

type Consumer interface {
	HandleFunc(typ HandleMessageType, fn HandleMessageFunc)
	Listen(ctx context.Context, opts *ConsumerOptions) error
}

type HandleMessageType string
type HandleMessageFunc func(ctx context.Context, data []byte) error

type consumer struct {
	channel  *amqp.Channel
	handlers map[HandleMessageType]HandleMessageFunc
}

func newConsumer(ch *amqp.Channel) Consumer {
	return &consumer{
		channel:  ch,
		handlers: make(map[HandleMessageType]HandleMessageFunc),
	}
}

func (c *consumer) HandleFunc(typ HandleMessageType, fn HandleMessageFunc) {
	c.handlers[typ] = fn
	log.Println("handler registered to type: ", typ)
}

func (c *consumer) Listen(ctx context.Context, opts *ConsumerOptions) error {

	messages, err := c.channel.Consume(
		opts.QueueName,
		opts.Consumer,
		opts.AutoAck,
		opts.Exclusive,
		opts.NoLocal,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}

	for message := range messages {

		var msg Message
		err := json.Unmarshal(message.Body, &msg)
		if err != nil {
			return err
		}

		handler, ok := c.handlers[msg.Type]
		if !ok {
			log.Printf("unhandled message: Type=%s\n", msg.Type)
		} else {
			err = handler(ctx, msg.Data)
			if err != nil {
				log.Println("error on handle message:", err.Error())
			}
		}

		if !opts.AutoAck {
			err = message.Ack(false)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
