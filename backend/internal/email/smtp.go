package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func (c *SMTPConfig) IsConfigured() bool {
	return c.Host != "" && c.Username != "" && c.Password != "" && c.From != ""
}

func SendOTP(cfg SMTPConfig, to string, otp string) error {
	subject := "RetailOS Super Admin Login OTP"
	body := fmt.Sprintf(
		"Your one-time password for RetailOS Super Admin login is:\n\n%s\n\nThis code expires in 5 minutes. If you did not request this, ignore this email.",
		otp,
	)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		cfg.From, to, subject, body,
	)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	// Port 465 = implicit TLS (SSL), Port 587 = STARTTLS
	if cfg.Port == "465" {
		return sendWithTLS(cfg, addr, to, []byte(msg))
	}
	return sendWithSTARTTLS(cfg, addr, to, []byte(msg))
}

// sendWithTLS connects over implicit TLS (port 465).
// Go's smtp.SendMail does not support this — it only does STARTTLS.
func sendWithTLS(cfg SMTPConfig, addr string, to string, msg []byte) error {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}

	tlsConn := tls.Client(conn, &tls.Config{ServerName: cfg.Host})
	if err := tlsConn.Handshake(); err != nil {
		conn.Close()
		return fmt.Errorf("tls handshake: %w", err)
	}

	client, err := smtp.NewClient(tlsConn, cfg.Host)
	if err != nil {
		tlsConn.Close()
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := client.Mail(cfg.Username); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return client.Quit()
}

// sendWithSTARTTLS uses port 587 with STARTTLS upgrade (what smtp.SendMail does,
// but with a connection timeout so it doesn't hang).
func sendWithSTARTTLS(cfg SMTPConfig, addr string, to string, msg []byte) error {
	conn, err := net.DialTimeout("tcp", addr, 10*time.Second)
	if err != nil {
		return fmt.Errorf("dial %s: %w", addr, err)
	}

	client, err := smtp.NewClient(conn, cfg.Host)
	if err != nil {
		conn.Close()
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if err := client.StartTLS(&tls.Config{ServerName: cfg.Host}); err != nil {
		return fmt.Errorf("starttls: %w", err)
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("smtp auth: %w", err)
	}

	if err := client.Mail(cfg.Username); err != nil {
		return fmt.Errorf("smtp MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT TO: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close data: %w", err)
	}

	return client.Quit()
}
