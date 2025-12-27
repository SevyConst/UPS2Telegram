package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type RequestBody struct {
	ChatID int64 `json:"chat_id"`
	Text string `json:"text"`
}

type ResponseBody struct {
	OK bool `json:"ok"`
	Error string `json:"description,omitempty"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result,omitempty"`
}

func sendToChatID(token string, chatID int64, message string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	requestBody := RequestBody {
		ChatID: chatID,
		Text: message,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("Can't marshal request body: %w", err)
	}

	client := &http.Client{ 
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Can't create http-request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Can't send http-request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Can't read answer: %w", err)
	}

	var responseBody ResponseBody 
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return fmt.Errorf("Can't parse answer: %w", err)
	}

	if !responseBody.OK {
		return fmt.Errorf("Telegram API error: %s", responseBody.Error)
	}

	log.Printf("Message '%s' has been sent to chat id %d", message, chatID)

	return nil
}

func SendToMultipleChats(token string, chatIDs []int64, message string) error {
	for _, chatID := range chatIDs {
		err := sendToChatID(token, chatID, message)
		if err != nil {
			return fmt.Errorf("Can't send to chat id %d: %w", chatID, err)
		}

		time.Sleep(100 * time.Millisecond)
	}
	return nil
}