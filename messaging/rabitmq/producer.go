package rabitmq

import (
	"context"
	"fmt"

	"github.com/rabbitmq/amqp091-go"

	"github.com/MiftahSalam/fa-fcs-go-lib/errors"
	"github.com/MiftahSalam/fa-fcs-go-lib/messaging"
)

type producerRMQ struct {
	common
	option *ProducerOptions
}

var (
	connectionPoolProducer = make(map[string]*producerRMQ)
)

func NewRMQProducer(opt *ProducerOptions) messaging.Producer {
	c, ok := connectionPoolProducer[opt.Name]
	if ok {
		return c
	}

	c = &producerRMQ{
		option: opt,
		common: common{
			err: make(chan error),
		},
	}

	connectionPoolProducer[opt.Name] = c

	return c
}

func (p *producerRMQ) SendMessage(ctx context.Context, exchange, topic string, message messaging.Message) error {
	select {
	case err := <-p.err:
		if err != nil {
			p.Reconnect()
		}
	default:
	}

	if p.conn == nil {
		return p.Reconnect()
	}

	msg := amqp091.Publishing{
		Headers:       amqp091.Table{"type": message.Body.Type},
		ContentType:   message.ContentType,
		CorrelationId: message.CorrelationID,
		Body:          message.Body.Data,
		Expiration:    message.Expired,
		Timestamp:     message.Timestamp,
	}
	if err := p.channel.PublishWithContext(ctx, exchange, topic, false, false, msg); err != nil {
		if err == amqp091.ErrClosed {
			p.Reconnect()
		}
		connErr := errors.ExtractError(errors.ErrConnection)
		return errors.New(connErr.HttpCode, connErr.Code, err.Error())
	}

	return nil
}

func (p *producerRMQ) BindExchange(dest, src, topic string) error {
	select {
	case err := <-p.err:
		if err != nil {
			p.Reconnect()
		}
	default:
	}

	if p.conn == nil {
		return p.Reconnect()
	}

	err := p.channel.ExchangeDeclare(
		dest,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("failed to declare destination exchange %s key %s with error: %v \n", dest, topic, err)

		connErr := errors.ExtractError(errors.ErrConnection)
		return errors.New(connErr.HttpCode, connErr.Code, err.Error())
	}

	err = p.channel.ExchangeDeclare(
		src,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("failed to declare source exchange %s key %s with error: %v \n", dest, topic, err)

		connErr := errors.ExtractError(errors.ErrConnection)
		return errors.New(connErr.HttpCode, connErr.Code, err.Error())
	}

	err = p.channel.ExchangeBind(dest, topic, src, false, nil)
	if err != nil {
		fmt.Printf("failed to bind exchange %s and %s key %s with error: %v \n", dest, src, topic, err)
		if err := p.Connect(); err != nil {
			return err
		}
	}

	return nil
}

func (p *producerRMQ) UnbindExchange(dest, src, topic string) error {
	select {
	case err := <-p.err:
		if err != nil {
			p.Reconnect()
		}
	default:
	}

	if p.conn == nil {
		return p.Reconnect()
	}

	err := p.channel.ExchangeUnbind(dest, topic, src, false, nil)
	if err != nil {
		fmt.Printf("failed to unbind exchange %s and %s key %s with error: %v \n", dest, src, topic, err)
		if err := p.Connect(); err != nil {
			return err
		}
	}

	return nil
}

func (p *producerRMQ) Connect() error {
	var err error

	p.conn, err = amqp091.Dial(p.option.Address)
	if err != nil {
		return err
	}

	go func() {
		<-p.conn.NotifyClose(make(chan *amqp091.Error))
		p.err <- errors.ErrConnection
	}()

	p.channel, err = p.conn.Channel()
	if err != nil {
		connErr := errors.ExtractError(errors.ErrConnection)
		return errors.New(connErr.HttpCode, connErr.Code, err.Error())
	}

	return nil
}

func (p *producerRMQ) Reconnect() error {
	if err := p.Connect(); err != nil {
		return err
	}

	return nil
}
