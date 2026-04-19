package port

import "context"

type Message struct {
	To      string
	Subject string
	Body    string
}

type Email interface {
	Send(ctx context.Context, message Message) error
}
