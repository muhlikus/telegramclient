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

	if cfg.Token == "" {
		return nil, errEmptyToken
	}

	// пока сюда вставил значения по умолчанию
	if cfg.BotApiScheme == "" {
		cfg.BotApiScheme = "https"
	}
	if cfg.BotApiHost == "" {
		cfg.BotApiHost = "api.telegram.org"
	}

	cfg.botApiPath = fmt.Sprintf("/bot%s", cfg.Token)
	cfg.HttpTimeout = 2000
	cfg.HttpTLSHandshakeTimeout = 500

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
