package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"sync"

	"github.com/rekib0023/event-horizon-mail-server/internal/config"
	"github.com/rekib0023/event-horizon-mail-server/internal/logger"
)

type Email struct {
	EmailType  string
	Recipients []string
	Subject    string
	Data       interface{}
}

type EmailService struct {
	conf *config.Config
	logg *logger.Logger
}

func NewEmailService(conf *config.Config, logg *logger.Logger) *EmailService {
	return &EmailService{
		conf: conf,
		logg: logg,
	}
}

func (e *EmailService) SendEmail(to, subject, templateName string, data interface{}) error {
	t, err := template.ParseFiles(fmt.Sprintf("internal/templates/%s", templateName))

	if err != nil {
		e.logg.Error("Error parsing template:", err)
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		e.logg.Error("Error executing template:", err)
		return err
	}

	body := buf.String()

	from := e.conf.SMTPUser
	// password := e.conf.SMTPPassword
	smtpHost := e.conf.SMTPServer
	smtpPort := e.conf.SMTPPort

	// auth := smtp.PlainAuth("", from, password, smtpHost)
	var auth smtp.Auth

	msg := []byte("Subject: " + subject + "\r\n" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
	if err != nil {
		e.logg.Error("Error sending email:", err)
		return err
	}

	e.logg.Info("Email sent successfully to", to)
	return nil
}

func (e *EmailService) SendBulkEmail(to []string, subject, templateName string, data interface{}) {
	var wg sync.WaitGroup
	for _, recipient := range to {
		wg.Add(1)
		go func(recipient string) {
			defer wg.Done()
			if err := e.SendEmail(recipient, subject, templateName, data); err != nil {
				e.logg.Error("Error sending email to", recipient, ":", err)
			}
		}(recipient)
	}
	wg.Wait()
}
