package main

import (
	"context"

	"github.com/MiftahSalam/fa-fcs-go-lib/messaging"
	"github.com/MiftahSalam/fa-fcs-go-lib/messaging/rabitmq"
)

func main() {
	options := &rabitmq.ProducerOptions{
		Address:  "amqp://guest:guest@localhost:5672/",
		Name:     "test-fa-fcs",
		Exchange: "test-exchange-fa-fcs",
		Routing:  "test.fa.fcs.*",
	}

	p := rabitmq.NewRMQProducer(options)

	err := p.Connect()
	if err != nil {
		panic("cannot connect to broker server")
	}

	m := messaging.Message{
		Queue: "speed-data",
		Body: messaging.MessageBody{
			Data: []byte("{\"sog\": 2.3}"),
			Type: "json",
		},
		ContentType: "application/json",
	}

	err = p.SendMessage(context.Background(), options.Routing, m)
	if err != nil {
		panic("cannot connect to broker server")
	}
}