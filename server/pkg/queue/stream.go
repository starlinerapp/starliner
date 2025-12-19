package queue

import (
	"errors"
	"github.com/nats-io/nats.go"
	"log"
)

type Stream string

const (
	Projects Stream = "projects"
)

func EnsureStream(js nats.JetStreamContext, name Stream, subjects []Subject) error {
	_, err := js.StreamInfo(string(name))
	if err != nil {
		s := make([]string, len(subjects))
		for i, subject := range subjects {
			s[i] = string(subject)
		}

		if errors.Is(err, nats.ErrStreamNotFound) {
			_, err := js.AddStream(&nats.StreamConfig{
				Name:      string(name),
				Subjects:  s,
				Storage:   nats.FileStorage,
				Retention: nats.LimitsPolicy,
			})
			if err != nil {
				return err
			}
			log.Printf("Created JetStream stream: %s for subjects: %v", name, subjects)
			return nil
		}
		return err
	}
	return nil
}
