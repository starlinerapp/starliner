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
	streamInfo, err := js.StreamInfo(string(name))
	if errors.Is(err, nats.ErrStreamNotFound) {
		// Stream doesn't exist, create it
		s := make([]string, len(subjects))
		for i, subject := range subjects {
			s[i] = string(subject)
		}

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
	} else if err != nil {
		return err
	}

	// Stream exists, check for missing subjects
	existingSubjects := make(map[string]struct{}, len(streamInfo.Config.Subjects))
	for _, subj := range streamInfo.Config.Subjects {
		existingSubjects[subj] = struct{}{}
	}

	var newSubjects []string
	for _, subj := range subjects {
		if _, ok := existingSubjects[string(subj)]; !ok {
			newSubjects = append(newSubjects, string(subj))
		}
	}

	if len(newSubjects) > 0 {
		// Update stream with new subjects
		updatedSubjects := append(streamInfo.Config.Subjects, newSubjects...)
		_, err := js.UpdateStream(&nats.StreamConfig{
			Name:      string(name),
			Subjects:  updatedSubjects,
			Storage:   streamInfo.Config.Storage,
			Retention: streamInfo.Config.Retention,
		})
		if err != nil {
			return err
		}
		log.Printf("Added new subjects %v to JetStream stream: %s", newSubjects, name)
	}

	return nil
}
