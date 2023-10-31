package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rekib0023/event-horizon-mail-server/internal/config"
	"github.com/rekib0023/event-horizon-mail-server/internal/email"
	"github.com/rekib0023/event-horizon-mail-server/internal/logger"
	"github.com/rekib0023/event-horizon-mail-server/internal/nats"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	conf := config.NewConfig()
	logg := logger.NewLogger()
	emailService := email.NewEmailService(conf, logg)
	natsSubscriber := nats.NewSubscriber(conf, logg, emailService)

	logg.Info("Email server is ready!")

	go natsSubscriber.SubscribeEmailSend()
	go natsSubscriber.SubscribeBulkEmailSend()

	select {}
}
