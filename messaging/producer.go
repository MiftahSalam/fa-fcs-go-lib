package messaging

import (
	"context"
)

type Producer interface {
	Connection
	SendMessage(ctx context.Context, exchange, topic string, message Message) error
	BindExchange(dest, src, topic string) error
	UnbindExchange(dest, src, topic string) error
}
