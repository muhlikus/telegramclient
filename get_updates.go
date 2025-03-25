package telegramclient

import (
	"encoding/json"
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
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending request: %w", err)
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

	var updates []Update
	err = json.Unmarshal(response.Result, &updates)
	if err != nil {
		return nil, fmt.Errorf("parsing updates JSON: %d", err)
	}

	return updates, nil
}
