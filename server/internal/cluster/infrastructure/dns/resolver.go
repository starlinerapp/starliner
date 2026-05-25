package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"starliner.app/internal/cluster/domain/port"
)

const pollInterval = 5 * time.Second

type Resolver struct{}

func NewResolver() port.DNS {
	return &Resolver{}
}

func (r *Resolver) WaitForHost(ctx context.Context, host string, expectedIP string) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		ips, err := net.LookupHost(host)
		if err == nil {
			for _, ip := range ips {
				if ip == expectedIP {
					return nil
				}
			}
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for %s to resolve to %s: %w", host, expectedIP, ctx.Err())
		case <-ticker.C:
		}
	}
}

func (r *Resolver) AllHostsResolve(hosts []string, expectedIP string) bool {
	for _, host := range hosts {
		ips, err := net.LookupHost(host)
		if err != nil {
			return false
		}

		found := false
		for _, ip := range ips {
			if ip == expectedIP {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}
