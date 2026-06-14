package value

import "time"

type CreateRunnerResult struct {
	Id        int64
	Token     string
	ExpiresAt time.Time
}
