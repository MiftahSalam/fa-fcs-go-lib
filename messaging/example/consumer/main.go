package main

import (
	"context"
	"fmt"

	"github.com/MiftahSalam/fa-fcs-go-lib/messaging"
	"github.com/MiftahSalam/fa-fcs-go-lib/messaging/rabitmq"
)

func main() {
	options := &rabitmq.ConsumerOptions{
		Address:  "amqp://guest:guest@localhost:5672/",
		Name:     "test-fa-fcs",
		Exchange: "test-exchange-fa-fcs",
		Routing:  "test.fa.fcs.*",
		Queues:   []string{"manager-osd"},
	}

	c := rabitmq.NewRMQConsumer(options)

	err := c.Connect()
	if err != nil {
		panic("cannot connect to broker server")
	}

	err = c.BindQueue()
	if err != nil {
		panic("cannot connect to broker server")
	}

	c.Start(context.Background(), options.Routing, func(msg messaging.Message) {
		fmt.Printf("receiver data '%s' from queue '%s'", msg.Body.Data, msg.Queue)
	})
}
