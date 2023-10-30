package messaging

type Message struct {
	Queue         string      `json:"queue"`
	ReplyTo       string      `json:"reply_to"`
	ContentType   string      `json:"content_type"`
	CorrelationID string      `json:"correlation_id"`
	Priority      uint8       `json:"priority"`
	Body          MessageBody `json:"body"`
}

type MessageBody struct {
	Data []byte `json:"data"`
	Type string `json:"type"`
}
