package value

import (
	"errors"
	"time"
)

var ErrInvalidRunnerRegistrationToken = errors.New("invalid runner registration token")

type CreateRunnerResult struct {
	Id        int64
	Token     string
	ExpiresAt time.Time
}
