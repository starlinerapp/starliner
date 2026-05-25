package dns

import (
	"context"
	"fmt"
	"net"
	"time"

	"starliner.app/internal/cluster/domain/port"
)

const (
	pollInterval   = 5 * time.Second
	publicResolver = "1.1.1.1:53"
	lookupTimeout  = 5 * time.Second
)

var publicDNSResolver = &net.Resolver{
	PreferGo: true,
	Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
		d := net.Dialer{Timeout: lookupTimeout}
		return d.DialContext(ctx, "udp", publicResolver)
	},
}

type Resolver struct{}

func NewResolver() port.DNS {
	return &Resolver{}
}

func (r *Resolver) WaitForHost(ctx context.Context, host string, expectedIP string) error {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		ips, err := lookupHost(ctx, host)
		if err == nil && includesIP(ips, expectedIP) {
			return nil
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
		ips, err := lookupHost(context.Background(), host)
		if err != nil || !includesIP(ips, expectedIP) {
			return false
		}
	}

	return true
}

func lookupHost(ctx context.Context, host string) ([]string, error) {
	lookupCtx, cancel := context.WithTimeout(ctx, lookupTimeout)
	defer cancel()

	return publicDNSResolver.LookupHost(lookupCtx, host)
}

func includesIP(ips []string, expectedIP string) bool {
	for _, ip := range ips {
		if ip == expectedIP {
			return true
		}
	}
	return false
}
