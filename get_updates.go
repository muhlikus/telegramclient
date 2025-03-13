package telegramclient

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) GetUpdates() ([]Update, error) {

	reqURL := url.URL{
		Scheme: c.cfg.botApiScheme,
		Host:   c.cfg.botApiHost,
		Path:   path.Join(c.cfg.botApiPath, getUpdatesMethod),
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), http.NoBody)
	if err != nil {
		return nil, err
	}

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
		return nil, err
	}

	if !response.OK {
		return nil, errors.New(response.Description)
	}

	var updates []Update
	err = json.Unmarshal(response.Result, &updates)
	if err != nil {
		return nil, err
	}
	return updates, nil
}
