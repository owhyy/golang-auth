package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFrom     string

	BaseURL string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	godotenv.Load()

	cfg.SMTPHost = os.Getenv("SMTP_HOST")
	cfg.SMTPPort = os.Getenv("SMTP_PORT")
	cfg.SMTPUsername = os.Getenv("SMTP_USERNAME")
	cfg.SMTPPassword = os.Getenv("SMTP_PASSWORD")
	cfg.SMTPFrom = os.Getenv("SMTP_FROM")
	cfg.BaseURL = os.Getenv("BASE_URL")

	if cfg.SMTPHost == "" || cfg.SMTPPort == "" || cfg.SMTPUsername == "" ||
		cfg.SMTPPassword == "" || cfg.SMTPFrom == "" || cfg.BaseURL == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return cfg, nil
}
