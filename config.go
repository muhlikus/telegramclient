package telegramclient

import "time"

type Config struct {
	Token                   string
	botApiScheme            string
	botApiHost              string
	botApiPath              string
	httpTimeout             time.Duration
	httpTLSHandshakeTimeout time.Duration
}
