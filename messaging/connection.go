package messaging

type Connection interface {
	Connect() error
	Reconnect() error
}
