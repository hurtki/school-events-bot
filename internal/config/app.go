package config

import (
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/DeanPDX/dotconfig"
)

type AppConfig struct {
	// ID of google spreadsheets document with schedule
	SpreadsheetsDocumentID string `env:"SCHEDULE_DOCUMENT_ID,required"`
	// Path where json schedule repository will store file
	JsonScheduleFileRepositoryPath string `env:"SCHEDULE_FILE_REPOSIOTORY_PATH,required"`
	// Path where json schedule repository will store file
	JsonPinnedMessageStateFileRepositoryPath string `env:"PINNED_MESSAGE_STATE_FILE_REPOSIOTORY_PATH,required"`

	// Interval for schedule poller
	SchedulePollerIntervalStr string `env:"SCHEDULE_POLLER_INTERVAL,required"`
	SchedulePollerInterval    time.Duration
}

func LoadAppConfig(src LoadSource) (AppConfig, error) {
	var cfg AppConfig
	var err error
	switch src {
	case EnviromentVariablesSource:
		var buf bytes.Buffer

		for _, e := range os.Environ() {
			// e = "KEY=VALUE"
			buf.WriteString(e)
			buf.WriteByte('\n')
		}

		cfg, err = dotconfig.FromReader[AppConfig](&buf)
	case EnvFileSource:
		cfg, err = dotconfig.FromFileName[AppConfig](".env")
	default:
		panic("wrong load source option")
	}

	duration, err := time.ParseDuration(cfg.SchedulePollerIntervalStr)
	if err != nil {
		return cfg, fmt.Errorf("invalid duration: %w", err)
	}
	cfg.SchedulePollerInterval = duration

	return cfg, err
}
