package rabitmq

import "github.com/rabbitmq/amqp091-go"

type CommonOptions struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}

type ProducerOptions struct {
	CommonOptions
}

type ConsumerOptions struct {
	CommonOptions
	ExchangeType string   `json:"exchange_type"`
	Routing      string   `json:"routing"`
	Exchange     string   `json:"exchange"`
	Queues       []string `json:"queues"`
}

type common struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	err     chan error
}
