package port

import "context"

type DNS interface {
	WaitForHost(ctx context.Context, host string, expectedIP string) error
	AllHostsResolve(hosts []string, expectedIP string) bool
}
