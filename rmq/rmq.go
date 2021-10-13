package rmq

import (
	"github.com/streadway/amqp"
	"go-subscriptions-workflow/util"
)

type Connection interface {
	ExchangeDeclare(opts *ExchangeOptions) error
	QueueDeclare(opts *QueueOptions) error
	NewConsumer() Consumer
	NewProducer() Producer
	Close()
}

type connection struct {
	amqpConn    *amqp.Connection
	amqpChannel *amqp.Channel
}

func New(cfg Config) Connection {
	conn, err := amqp.Dial(cfg.URL())
	util.PanicOnError(err)
	ch, err := conn.Channel()
	util.PanicOnError(err)
	return &connection{
		amqpConn:    conn,
		amqpChannel: ch,
	}
}

func (c *connection) Close() {
	util.HandleClose(c.amqpChannel)
	util.HandleClose(c.amqpConn)
}

func (c *connection) ExchangeDeclare(opts *ExchangeOptions) error {
	return c.amqpChannel.ExchangeDeclare(
		opts.Name,
		opts.Kind,
		opts.Durable,
		opts.AutoDelete,
		opts.Internal,
		opts.NoWait,
		opts.Args,
	)
}

func (c *connection) QueueDeclare(opts *QueueOptions) error {
	queue, err := c.amqpChannel.QueueDeclare(
		opts.Name,
		opts.Durable,
		opts.AutoDelete,
		opts.Exclusive,
		opts.NoWait,
		opts.Args,
	)
	if err != nil {
		return err
	}

	if opts.BindOptions != nil {

		bindOpts := opts.BindOptions

		err = c.amqpChannel.QueueBind(
			queue.Name,
			bindOpts.RoutingKey,
			bindOpts.ExchangeName,
			bindOpts.NoWait,
			bindOpts.Args,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *connection) NewConsumer() Consumer {
	return newConsumer(c.amqpChannel)
}

func (c *connection) NewProducer() Producer {
	return newProducer(c.amqpChannel)
}
