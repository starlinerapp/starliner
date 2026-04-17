package email

import (
	"context"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
)

type Client struct {
	cfg *conf.Config
}

func NewClient(cfg *conf.Config) port.Email {
	return &Client{cfg: cfg}
}

func (c *Client) Send(ctx context.Context, message port.Message) error {
	recipient, err := mail.ParseAddress(message.To)
	if err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}

	client, err := smtp.Dial(c.cfg.SmtpHost + ":" + c.cfg.SmtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Quit()

	if c.cfg.SmtpUsername != "" && c.cfg.SmtpPassword != "" {
		auth := smtp.PlainAuth("", c.cfg.SmtpUsername, c.cfg.SmtpPassword, c.cfg.SmtpHost)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	if err = client.Mail(c.cfg.FromMail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(recipient.Address); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data: %w", err)
	}

	_, err = fmt.Fprintf(wc,
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		c.cfg.FromMail, recipient.Address, message.Subject,
	)
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}

	if _, err = io.WriteString(wc, message.Body); err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	if err = wc.Close(); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}

	return nil
}
