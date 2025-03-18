package telegramclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) SendMessage(chatID int, text string) (*Message, error) {
	newMsg := struct {
		ChatID int    `json:"chat_id"`
		Text   string `json:"text"`
	}{
		ChatID: chatID,
		Text:   text,
	}

	reqURL := url.URL{
		Scheme: c.cfg.botApiScheme,
		Host:   c.cfg.botApiHost,
		Path:   path.Join(c.cfg.botApiPath, sendMessageMethod),
	}

	newMsgJSON, err := json.Marshal(newMsg)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, reqURL.String(), bytes.NewReader(newMsgJSON))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("parsing response JSON: %d", err)
	}

	if !response.OK {
		return nil, fmt.Errorf("response not OK: %s", response.Description)
	}

	var message Message
	err = json.Unmarshal(response.Result, &message)
	if err != nil {
		return nil, fmt.Errorf("parsing message JSON: %d", err)
	}
	return &message, nil
}
