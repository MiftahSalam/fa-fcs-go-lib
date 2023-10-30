package rabitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"

	"github.com/MiftahSalam/fa-fcs-go-lib/errors"
	"github.com/MiftahSalam/fa-fcs-go-lib/messaging"
)

type consumerRMQ struct {
	common
	option *ConsumerOptions
}

var (
	connectionPoolConsumer = make(map[string]*consumerRMQ)
)

func NewRMQConsumer(opt *ConsumerOptions) messaging.Consumer {
	c, ok := connectionPoolConsumer[opt.Name]
	if ok {
		return c
	}

	c = &consumerRMQ{
		option: opt,
		common: common{
			err: make(chan error),
		},
	}

	connectionPoolConsumer[opt.Name] = c

	return c
}

func (consumer *consumerRMQ) Start(ctx context.Context, topic string, handler messaging.ConsumeHandler) error {
	forever := make(chan bool)

	m, err := consumer.consume()
	if err != nil {
		return err
	}

	for q, d := range m {
		go consumer.handleConsumed(q, d, handler)
	}

	<-forever

	return nil
}

func (consumer *consumerRMQ) handleConsumed(q string, delivery <-chan amqp091.Delivery, handler messaging.ConsumeHandler) {
	fmt.Println("handleConsumed for queue: ", q)

	for {
		for d := range delivery {
			msg := messaging.Message{
				Queue: q,
				Body: messaging.MessageBody{
					Data: d.Body,
					Type: d.Headers["type"].(string),
				},
			}

			handler(msg)

			d.Ack(false)
		}

		if err := <-consumer.err; err != nil {
			fmt.Println("rabitmq connection broken with error: ", err.Error())

			err = consumer.Reconnect()
			if err == nil {
				fmt.Println("rabitmq regain connection")
				deliveries, err := consumer.consume()
				if err == nil {
					delivery = deliveries[q]
				}
			} else {
				go func() {
					consumer.err <- errors.ErrConnection
				}()
			}

			time.Sleep(1 * time.Second)
		}
	}

}

func (consumer *consumerRMQ) consume() (map[string]<-chan amqp091.Delivery, error) {
	m := make(map[string]<-chan amqp091.Delivery)

	for _, q := range consumer.option.Queues {
		deliveries, err := consumer.channel.Consume(q, "", false, false, false, false, nil)
		if err != nil {
			return nil, err
		}

		m[q] = deliveries
	}

	return m, nil
}

func (consumer *consumerRMQ) Connect() error {
	var err error

	consumer.conn, err = amqp091.Dial(consumer.option.Address)
	if err != nil {
		return err
	}

	go func() {
		<-consumer.conn.NotifyClose(make(chan *amqp091.Error))
		consumer.err <- errors.ErrConnection
	}()

	consumer.channel, err = consumer.conn.Channel()
	if err != nil {
		connErr := errors.ExtractError(errors.ErrConnection)
		return errors.New(connErr.HttpCode, connErr.Code, err.Error())
	}

	err = consumer.channel.ExchangeDeclare(
		consumer.option.Exchange,
		consumer.option.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		connErr := errors.ExtractError(errors.ErrConnection)
		return errors.New(connErr.HttpCode, connErr.Code, err.Error())
	}

	return nil
}

func (consumer *consumerRMQ) BindQueue() error {
	for _, q := range consumer.option.Queues {
		if _, err := consumer.channel.QueueDeclare(q, true, false, false, false, nil); err != nil {
			connErr := errors.ExtractError(errors.ErrConnection)
			return errors.New(connErr.HttpCode, connErr.Code, err.Error())
		}
		if err := consumer.channel.QueueBind(q, consumer.option.Routing, consumer.option.Exchange, false, nil); err != nil {
			connErr := errors.ExtractError(errors.ErrConnection)
			return errors.New(connErr.HttpCode, connErr.Code, err.Error())
		}
	}

	return nil
}

func (consumer *consumerRMQ) Reconnect() error {
	err := consumer.Connect()
	if err == nil {
		return consumer.BindQueue()
	}

	return err
}
