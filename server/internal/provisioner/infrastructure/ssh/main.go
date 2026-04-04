package ssh

import (
	"context"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"starliner.app/internal/provisioner/domain/port"
	"time"
)

type SSH struct{}

func NewSSH() port.SSH {
	return &SSH{}
}

func (s *SSH) WaitForSSH(ip string, timeout time.Duration) error {
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

func (s *SSH) OpenTTY(
	ctx context.Context,
	user string,
	ip string,
	pemKey []byte,
	terminalSize <-chan port.TerminalSize,
) (stdin io.WriteCloser, stdout io.ReadCloser, err error) {
	signer, err := ssh.ParsePrivateKey(pemKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	go func() {
		defer func() {
			_ = stdinReader.Close()
			_ = stdoutWriter.Close()
		}()

		var client *ssh.Client
		for {
			select {
			case <-ctx.Done():
				_ = stdoutWriter.CloseWithError(ctx.Err())
				return
			default:
			}

			c, err := ssh.Dial("tcp", net.JoinHostPort(ip, "22"), config)
			if err != nil {
				time.Sleep(2 * time.Second)
				continue
			}
			client = c
			break
		}
		defer func() { _ = client.Close() }()

		session, err := client.NewSession()
		if err != nil {
			_ = stdoutWriter.CloseWithError(fmt.Errorf("failed to create session: %w", err))
			return
		}
		defer func() { _ = session.Close() }()

		session.Stdin = stdinReader
		session.Stdout = stdoutWriter
		session.Stderr = stdoutWriter

		initialSize := port.TerminalSize{Columns: 80, Rows: 24}
		if err := session.RequestPty("xterm-256color", int(initialSize.Rows), int(initialSize.Columns), ssh.TerminalModes{
			ssh.ECHO:          1,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}); err != nil {
			_ = stdoutWriter.CloseWithError(fmt.Errorf("failed to request pty: %w", err))
			return
		}

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case size, ok := <-terminalSize:
					if !ok {
						return
					}
					_ = session.WindowChange(int(size.Rows), int(size.Columns))
				}
			}
		}()

		if err := session.Shell(); err != nil {
			_ = stdoutWriter.CloseWithError(fmt.Errorf("failed to start shell: %w", err))
			return
		}

		sessionDone := make(chan error, 1)
		go func() { sessionDone <- session.Wait() }()

		select {
		case <-ctx.Done():
			_ = stdoutWriter.CloseWithError(ctx.Err())
		case err := <-sessionDone:
			if err != nil {
				_ = stdoutWriter.CloseWithError(err)
			}
		}
	}()

	return stdinWriter, stdoutReader, nil
}
