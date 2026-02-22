package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds the application configuration
type Config struct {
	SlackWebhookURL      string
	SlackWebhookURLError string
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig loads configuration from environment variables
func LoadConfig() {
	// .envファイルから環境変数を読み込む。ファイルがなくてもエラーにはしない。
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig.SlackWebhookURL = os.Getenv("SLACK_WEBHOOK_URL")
	AppConfig.SlackWebhookURLError = os.Getenv("SLACK_WEBHOOK_URL_FOR_ERROR")

	if AppConfig.SlackWebhookURL == "" || AppConfig.SlackWebhookURLError == "" {
		log.Fatal("SLACK_WEBHOOK_URL and SLACK_WEBHOOK_URL_FOR_ERROR must be set")
	}
}
