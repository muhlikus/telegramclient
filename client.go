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
	if cfg.botApiScheme == "" {
		cfg.botApiScheme = "https"
	}
	if cfg.botApiHost == "" {
		cfg.botApiHost = "api.telegram.org"
	}

	cfg.botApiPath = fmt.Sprintf("/bot%s", cfg.Token)
	cfg.httpTimeout = 2000
	cfg.httpTLSHandshakeTimeout = 500

	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout: cfg.httpTLSHandshakeTimeout * time.Millisecond,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: cfg.httpTimeout * time.Millisecond,
		},
		cfg: cfg,
	}, nil
}
