package email

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/mail"
	"net/smtp"
	"starliner.app/internal/api/conf"
	"starliner.app/internal/api/domain/port"
	"strings"
)

const starlinerLogoURL = "https://starliner-596451156994-eu-north-1-an.s3.eu-north-1.amazonaws.com/starliner-logo.png"

type Client struct {
	cfg      *conf.Config
	renderer *Renderer
}

func NewClient(cfg *conf.Config, renderer *Renderer) port.Email {
	return &Client{
		cfg:      cfg,
		renderer: renderer,
	}
}

type renderInviteData struct {
	OrganizationName string
	InviteLink       string
	LogoURL          string
}

func (c *Client) SendInvite(to string, inviteData port.InviteData) error {
	html, err := c.renderer.Render("invite.html", renderInviteData{
		OrganizationName: inviteData.OrganizationName,
		InviteLink:       inviteData.InviteLink,
		LogoURL:          starlinerLogoURL,
	})
	if err != nil {
		return err
	}

	return c.send(&message{
		To:      to,
		Subject: fmt.Sprintf("Join %s on Starliner", inviteData.OrganizationName),
		Body:    html,
	})
}

type renderVerifyData struct {
	VerificationLink string
	LogoURL          string
}

func (c *Client) SendVerificationEmail(to string, verifyData port.VerifyData) error {
	html, err := c.renderer.Render("verify.html", renderVerifyData{
		VerificationLink: verifyData.VerificationLink,
		LogoURL:          starlinerLogoURL,
	})
	if err != nil {
		return err
	}

	return c.send(&message{
		To:      to,
		Subject: "[Starliner] Verify your email",
		Body:    html,
	})
}

type renderResetData struct {
	PasswordResetLink string
	LogoURL           string
}

func (c *Client) SendResetPassword(to string, resetData port.ResetData) error {
	html, err := c.renderer.Render("reset.html", renderResetData{
		PasswordResetLink: resetData.PasswordResetLink,
		LogoURL:           starlinerLogoURL,
	})
	if err != nil {
		return err
	}

	return c.send(&message{
		To:      to,
		Subject: "[Starliner] Reset your password",
		Body:    html,
	})
}

type message struct {
	To      string
	Subject string
	Body    string
}

func (c *Client) send(message *message) error {
	recipient, err := mail.ParseAddress(message.To)
	if err != nil {
		return fmt.Errorf("invalid recipient address: %w", err)
	}

	client, err := smtp.Dial(c.cfg.SmtpHost + ":" + c.cfg.SmtpPort)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer func() {
		if err != nil {
			_ = client.Close()
		} else {
			_ = client.Quit()
		}
	}()

	if c.cfg.SmtpTLSEnabled {
		if err = client.StartTLS(&tls.Config{ServerName: c.cfg.SmtpHost}); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
		auth := smtp.PlainAuth("", c.cfg.SmtpUsername, c.cfg.SmtpPassword, c.cfg.SmtpHost)
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	if err = client.Mail(c.cfg.SenderMail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(recipient.Address); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	wc, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data: %w", err)
	}

	subject := sanitizeHeaderSubject(message.Subject)
	_, err = fmt.Fprintf(wc,
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n",
		c.cfg.SenderMail, recipient.Address, subject,
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

// sanitizeHeaderSubject replaces CR/LF in values interpolated into raw SMTP
// header lines so they cannot inject additional headers.
func sanitizeHeaderSubject(s string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(s)
}
