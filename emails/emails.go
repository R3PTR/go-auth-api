// emails.go
package emails

import (
	"gopkg.in/gomail.v2"
)

// EmailSender represents an email sender.
type EmailSender struct {
	From     string
	Host     string
	Port     int
	Username string
	Password string
}

// NewEmailSender creates a new EmailSender instance.
func NewEmailSender(from, host string, port int, username, password string) *EmailSender {
	return &EmailSender{
		From:     from,
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}
}

// SendEmail sends an email.
func (e *EmailSender) SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(e.Host, e.Port, e.Username, e.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
