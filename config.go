package telegramclient

import "time"

type Config struct {
	Token                   string `env:"TELEGRAM_BOT_TOKEN,notEmpty"`
	BotApiScheme            string `env:"TELEGRAM_BOT_API_SCHEME" envDefault:"https"`
	BotApiHost              string `env:"TELEGRAM_BOT_API_HOST" envDefault:"api.telegram.org"`
	botApiPath              string
	HttpTimeout             time.Duration `env:"TELEGRAM_BOT_HTTP_TIMEOUT" envDefault:"2s"`
	HttpTLSHandshakeTimeout time.Duration `env:"TELEGRAM_BOT_TLS_TIMEOUT" envDefault:"500ms"`
}
