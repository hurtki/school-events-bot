package config

import (
	"bytes"
	"os"

	"github.com/DeanPDX/dotconfig"
)

type LoadSource uint8

const (
	EnviromentVariablesSource LoadSource = iota
	EnvFileSource
)

type BotConfig struct {
	// bot that app will use to operate channel
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN,required"`
	// channel, where bot will send updates
	UpdatesChannel int64 `env:"UPDATES_TELEGRAM_CHANNEL_ID,required"`
}

func LoadBotConfig(src LoadSource) (BotConfig, error) {
	switch src {
	case EnviromentVariablesSource:
		var buf bytes.Buffer

		for _, e := range os.Environ() {
			// e = "KEY=VALUE"
			buf.WriteString(e)
			buf.WriteByte('\n')
		}

		return dotconfig.FromReader[BotConfig](&buf)
	case EnvFileSource:
		return dotconfig.FromFileName[BotConfig](".env")
	default:
		panic("wrong load source option")
	}
}
