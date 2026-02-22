package main

import (
	"fmt"
	"log"
	"strings"
)

func main() {
	// 1. 設定を読み込む
	LoadConfig()
	log.Println("Configuration loaded.")

	// 2. スクレイピングを実行
	log.Println("Starting to scrape available slots...")
	availableSlots, err := ScrapeAvailableSlots()
	if err != nil {
		// 3. エラー処理
		errorMsg := fmt.Sprintf("Error during scraping: %v", err)
		log.Print(errorMsg)
		// エラー用Webhookに通知
		if sendErr := SendMessage(AppConfig.SlackWebhookURLError, errorMsg); sendErr != nil {
			log.Printf("Failed to send error message to Slack: %v", sendErr)
		}
		// エラーで終了する場合、ここでpanicやos.Exit(1)を呼ぶこともできる
		return
	}

	// 4. 空き枠が見つからなかった場合
	if len(availableSlots) == 0 {
		log.Println("No available slots found.")
		return
	}

	// 5. 空き枠が見つかった場合
	log.Printf("Found %d available slots!", len(availableSlots))
	var messageBuilder strings.Builder
	messageBuilder.WriteString("体操教室の空きが見つかりました！\n\n")

	for _, slot := range availableSlots {
		messageBuilder.WriteString(slot.String())
		messageBuilder.WriteString("\n")
	}
	messageBuilder.WriteString("予約はこちらから: " + baseURL)

	// Slackに通知
	log.Println("Sending notification to Slack...")
	if err := SendMessage(AppConfig.SlackWebhookURL, messageBuilder.String()); err != nil {
		errorMsg := fmt.Sprintf("Failed to send success message to Slack: %v", err)
		log.Println(errorMsg)
		// こちらもエラー用Webhookに通知した方が親切
		if sendErr := SendMessage(AppConfig.SlackWebhookURLError, errorMsg); sendErr != nil {
			log.Printf("Failed to send secondary error message to Slack: %v", sendErr)
		}
	} else {
		log.Println("Successfully sent notification to Slack.")
	}
}
