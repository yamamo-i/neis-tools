package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackMessage defines the structure for a Slack message payload.
type SlackMessage struct {
	Text string `json:"text"`
}

// SendMessage sends a message to the specified Slack webhook URL.
func SendMessage(webhookURL string, message string) error {
	if webhookURL == "" {
		return fmt.Errorf("slack webhook URL is not set")
	}

	payload := SlackMessage{
		Text: message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal slack message payload: %w", err)
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send message to slack: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-200 status code from slack: %s", resp.Status)
	}

	return nil
}
