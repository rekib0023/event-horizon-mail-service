package nats

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rekib0023/event-horizon-mail-server/internal/config"
	"github.com/rekib0023/event-horizon-mail-server/internal/email"
	"github.com/rekib0023/event-horizon-mail-server/internal/logger"
)

type Subscriber struct {
	conf         *config.Config
	logg         *logger.Logger
	emailService *email.EmailService
	nc           *nats.Conn
}

func NewSubscriber(conf *config.Config, logg *logger.Logger, emailService *email.EmailService) *Subscriber {
	nc, err := nats.Connect(conf.NATSURL)
	if err != nil {
		logg.Fatal("Error connecting to NATS server:", err)
	}

	return &Subscriber{
		conf:         conf,
		logg:         logg,
		emailService: emailService,
		nc:           nc,
	}
}

func (s *Subscriber) SubscribeEmailSend() error {
	return s.subscribe("email.send", false)
}

func (s *Subscriber) SubscribeBulkEmailSend() error {
	return s.subscribe("email.bulk_send", true)
}

func (s *Subscriber) subscribe(subject string, isBulk bool) error {
	_, err := s.nc.Subscribe(subject, func(msg *nats.Msg) {
		var email email.Email
		err := json.Unmarshal(msg.Data, &email)
		if err != nil {
			s.logg.Fatal(err)
		}

		var templateName string
		switch email.EmailType {
		case "confirmation":
			templateName = "confirmation.html"
		case "reminder":
			templateName = "reminder.html"
		}

		if isBulk {
			s.emailService.SendBulkEmail(email.Recipients, email.Subject, templateName, email.Data)
		} else {
			err = s.emailService.SendEmail(email.Recipients[0], email.Subject, templateName, email.Data)
			if err != nil {
				s.logg.Fatal(err)
			}
		}

		fmt.Println("Email sent successfully")
	})

	if err != nil {
		return err
	}

	return nil
}
