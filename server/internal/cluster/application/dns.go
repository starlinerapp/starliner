package application

import (
	"context"
	"fmt"
	"net"
	"time"
)

const dnsPollInterval = 5 * time.Second

func waitForDNS(ctx context.Context, host string, expectedIP string) error {
	ticker := time.NewTicker(dnsPollInterval)
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

func allHostsResolve(hosts []string, expectedIP string) bool {
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
