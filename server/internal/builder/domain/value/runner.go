package value

import "errors"

var (
	ErrInvalidRunnerCapacity = errors.New("invalid runner capacity")
	ErrRunnerDeleted         = errors.New("runner deleted")
)
