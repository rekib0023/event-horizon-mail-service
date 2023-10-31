package config

import (
	"os"
)

type Config struct {
	SMTPServer   string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	NATSURL      string
}

func NewConfig() *Config {
	return &Config{
		SMTPServer:   os.Getenv("SMTP_SERVER"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		NATSURL:      os.Getenv("NATSURL"),
	}
}
