package telegramclient

import "errors"

var (
	errEmptyToken        = errors.New("token is empty")
	errEmptyBotApiScheme = errors.New("bot api scheme is empty")
	errEmptyBotApiHost   = errors.New("bot api host is empty")
)
