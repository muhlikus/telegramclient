package telegramclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/caarlos0/env/v11"
)

type Client struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) (*Client, error) {
	err := env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	// if cfg.Token == "" {
	// 	return nil, errEmptyToken
	// }

	// // пока сюда вставил значения по умолчанию
	// if cfg.botApiScheme == "" {
	// 	cfg.botApiScheme = "https"
	// }

	// if cfg.botApiHost == "" {
	// 	cfg.botApiHost = "api.telegram.org"
	// }

	cfg.botApiPath = fmt.Sprintf("/bot%s", cfg.Token)
	// cfg.httpTimeout = 2000
	// cfg.httpTLSHandshakeTimeout = 500

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
