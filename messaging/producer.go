package messaging

import (
	"context"
)

type Producer interface {
	Connection
	SendMessage(ctx context.Context, topic string, message Message) error
}