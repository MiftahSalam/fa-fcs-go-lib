package main

import (
	"context"
	"fmt"
	"time"

	"github.com/MiftahSalam/fa-fcs-go-lib/messaging"
	"github.com/MiftahSalam/fa-fcs-go-lib/messaging/rabitmq"
)

func main() {
	options := &rabitmq.ConsumerOptions{
		CommonOptions: rabitmq.CommonOptions{
			Address: "amqp://guest:guest@localhost:5672/",
			Name:    "test-fa-fcs",
		},
		Exchange:     "test-exchange-fa-fcs",
		ExchangeType: "topic",
		Routing:      "test.fa.fcs.*",
		Queues:       []string{"manager-osd"},
	}

	c := rabitmq.NewRMQConsumer(options)

	err := fmt.Errorf("pending connection error")
	for {
		if err != nil {
			fmt.Println("Connecting to broker server...")
			err = c.Connect()
			time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Connected to broker server")
			break
		}
	}

	err = c.BindQueue()
	if err != nil {
		panic("cannot connect to broker server")
	}

	c.Start(context.Background(), options.Routing, func(msg messaging.Message) {
		fmt.Printf("receiver data '%s' from queue '%s'", msg.Body.Data, msg.Queue)
	})
}
