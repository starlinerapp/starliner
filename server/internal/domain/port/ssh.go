package port

import "time"

type SSH interface {
	WaitForSSH(ip string, timeout time.Duration) error
}
