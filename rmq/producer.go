package rmq

import (
	"github.com/streadway/amqp"
	"time"
)

type Producer interface{
	Send(opts *PublisherOptions, msg *Message) error
}

type producer struct {
	channel *amqp.Channel
}

func newProducer(ch *amqp.Channel) Producer {
	return &producer{
		channel: ch,
	}
}

func (p *producer) Send(opts *PublisherOptions, msg *Message) error {

	message := amqp.Publishing{
		ContentType: "application/json",
		Timestamp:   time.Now(),
		Body:        msg.Bytes(),
	}

	if opts.Persistent {
		message.DeliveryMode = amqp.Persistent
	}

	return p.channel.Publish(
		opts.ExchangeName,
		opts.RoutingKey,
		opts.Mandatory,
		opts.Immediate,
		message,
	)
}
