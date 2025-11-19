package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
)

type EmailConfig struct {
	EmailFrom     string
	EmailAlias    string
	EmailPassword string
	SMTPHost      string
	SMTPPort      string
}

type EmailClient struct {
	client *smtp.Client
	config EmailConfig
}

type Email struct {
	To       string
	Subject  string
	Body     string
	Template string
}

func NewClient(cfg EmailConfig) (EmailClient, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         cfg.SMTPHost,
	}

	conn, err := tls.Dial("tcp", cfg.SMTPHost+":"+cfg.SMTPPort, tlsConfig)
	if err != nil {
		return EmailClient{}, err
	}

	client, err := smtp.NewClient(conn, cfg.SMTPHost)
	if err != nil {
		return EmailClient{}, err
	}

	// authenticate client
	auth := smtp.PlainAuth("", cfg.EmailFrom, cfg.EmailPassword, cfg.SMTPHost)
	if err := client.Auth(auth); err != nil {
		return EmailClient{}, err
	}

	log.Printf("SMTP client for email: %s connected", cfg.EmailFrom)

	return EmailClient{
		client: client,
		config: cfg,
	}, nil
}

// if content type is text/html then the email body will be ignored
func (c *EmailClient) SendMail(param Email) error {
	// set sender
	if err := c.client.Mail(c.config.EmailFrom); err != nil {
		return err
	}

	// set recipient
	if err := c.client.Rcpt(param.To); err != nil {
		return err
	}

	// switch content type and body
	var contentType string
	var contentBody string

	if param.Template != "" {
		contentType = "text/html"
		contentBody = param.Template
	} else {
		contentType = "text/plain"
		contentBody = param.Body
	}

	// send data
	w, err := c.client.Data()
	if err != nil {
		return err
	}

	headers := c.buildHeaders(param, contentType)

	fullMessage := headers + contentBody

	if _, err := w.Write([]byte(fullMessage)); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	log.Printf("email sent from: %s to: %s", c.config.EmailFrom, param.To)

	return nil

}

func (c *EmailClient) buildHeaders(param Email, contentType string) string {
	headers := "From: \"" + c.config.EmailAlias + "\" <" + c.config.EmailFrom + ">\r\n" +
		"To: " + param.To + "\r\n" +
		"Subject: " + param.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: " + contentType + "; charset=\"UTF-8\"\r\n" +
		"\r\n"

	return headers
}

func (c *EmailClient) Quit() error {
	if err := c.client.Quit(); err != nil {
		return fmt.Errorf("unable to close SMTP connection, err: %s", err.Error())
	}

	return nil
}
