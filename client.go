package telegramclient

import (
	"crypto/tls"
	"errors"
	"net/http"
	"time"
)

const (
	httpTimeout             = 2000
	httpTLSHandshakeTimeout = 500
	queryTemplate           = "https://api.telegram.org/bot%s/%s" // https://api.telegram.org/bot<token>/METHOD_NAME
)

type Client struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) (*Client, error) {

	if cfg.Token == "" {
		return nil, errors.New("the token cannot be empty")
	}

	return &Client{
		client: &http.Client{
			Transport: &http.Transport{
				TLSHandshakeTimeout: httpTLSHandshakeTimeout * time.Millisecond,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			Timeout: httpTimeout * time.Millisecond,
		},
		cfg: cfg,
	}, nil
}
