package ssh

import (
	"fmt"
	"net"
	"starliner.app/internal/provisioner/domain/port"
	"time"
)

type SSH struct{}

func NewSSH() port.SSH {
	return &SSH{}
}

func (ssh *SSH) WaitForSSH(ip string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "22"), 5*time.Second)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for ssh on %s", ip)
		}
		time.Sleep(5 * time.Second)
	}
}
