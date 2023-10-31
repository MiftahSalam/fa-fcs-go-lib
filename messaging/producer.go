package messaging

import (
	"context"
)

type Producer interface {
	Connection
	SendMessage(ctx context.Context, exchange, topic string, message Message) error
}
