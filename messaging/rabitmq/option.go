package rabitmq

import "github.com/rabbitmq/amqp091-go"

type ProducerOptions struct {
	Address  string `json:"address"`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Routing  string `json:"routing"`
}

type ConsumerOptions struct {
	Address  string   `json:"address"`
	Name     string   `json:"name"`
	Exchange string   `json:"exchange"`
	Routing  string   `json:"routing"`
	Queues   []string `json:"queues"`
}

type common struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	err     chan error
}
