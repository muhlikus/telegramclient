package telegramclient

import (
	"encoding/json"
	"errors"
	"fmt"
)

func (c *Client) GetUpdates() ([]Update, error) {

	const op = "getUpdates"

	query := fmt.Sprintf(queryTemplate, c.cfg.Token, op)
	resp, err := c.client.Get(query)
	if err != nil {
		return []Update{}, err
	}
	defer resp.Body.Close()

	var response Response
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return []Update{}, err
	}

	if !response.OK {
		return []Update{}, errors.New(response.Description)
	}

	var updates []Update
	err = json.Unmarshal(response.Result, &updates)
	if err != nil {
		return []Update{}, err
	}
	return updates, nil
}
