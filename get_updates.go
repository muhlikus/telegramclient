package telegramclient

import (
	"encoding/json"
	"fmt"
)

func (c *Client) GetUpdates() ([]Update, error) {

	const op = "getUptates"

	query := fmt.Sprintf(queryTemplate, c.cfg.Token, op)
	resp, err := c.client.Get(query)
	if err != nil {
		return []Update{}, err
	}
	defer resp.Body.Close()

	var updates []Update
	err = json.NewDecoder(resp.Body).Decode(&updates)
	if err != nil {
		return []Update{}, err
	}

	return updates, nil
}
