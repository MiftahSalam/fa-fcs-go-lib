package messaging

import (
	"context"
)

type ConsumeHandler func(msg Message)

type Consumer interface {
	Connection
	Start(ctx context.Context, topic string, handler ConsumeHandler) error
	BindQueue() error
}
