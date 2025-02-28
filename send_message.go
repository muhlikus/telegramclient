package telegramclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Client) SendMessage(chatID int, text string) (Message, error) {
	const op = "sendMessage"

	type newMessage struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}

	newMsg := newMessage{
		ChatID: chatID,
		Text:   text,
	}

	query := fmt.Sprintf(queryTemplate, c.cfg.Token, op)
	newMsgJSON, err := json.Marshal(newMsg)
	if err != nil {
		return Message{}, err
	}

	resp, err := c.client.Post(query, "application/json", bytes.NewReader(newMsgJSON))
	if err != nil {
		return Message{}, err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return Message{}, err
	}

	if !response.OK {
		return Message{}, errors.New(response.Description)
	}

	var message Message
	err = json.Unmarshal(response.Result, &message)
	if err != nil {
		return Message{}, err
	}
	return message, nil
}
