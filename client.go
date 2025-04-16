package telegramclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) (*Client, error) {
	err := cfg.validate()

	if err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	cfg.botApiPath = fmt.Sprintf("/bot%s", cfg.Token)

	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout: cfg.HttpTLSHandshakeTimeout * time.Millisecond,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: cfg.HttpTimeout * time.Millisecond,
		},
		cfg: cfg,
	}, nil
}
