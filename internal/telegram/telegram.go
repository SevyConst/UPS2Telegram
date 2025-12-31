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
	} `json:"result"`
}

func sendToChatID(token string, chatID int64, message string) {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token)

	requestBody := RequestBody {
		ChatID: chatID,
		Text: message,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Printf("Chat id: %d - can't marshal request body: %v", chatID, err)
		return
	}

	client := &http.Client{ 
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Chat id: %d - can't create http-request: %v", chatID, err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Chat id: %d - can't send http-request: %v", chatID, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Chat id: %d - HTTP status %d", chatID, resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Chat id: %d - can't read answer: %v", chatID, err)
		return
	}

	var responseBody ResponseBody 
	if err := json.Unmarshal(body, &responseBody); err != nil {
		log.Printf("Chat id: %d - can't parse answer: %v", chatID, err)
		return
	}

	if !responseBody.OK {
		log.Printf("Chat id: %d - Telegram API error: %s", chatID, responseBody.Error)
		return
	}

	log.Printf("Chat id: %d - message '%s' has been sent to chat id %d", chatID, message, chatID)
}

func SendToMultipleChats(token string, chatIDs []int64, message string) error {
	for _, chatID := range chatIDs {
		sendToChatID(token, chatID, message)

		time.Sleep(100 * time.Millisecond)
	}
	return nil
}