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
	ScheduleWorkerIntervalStr string `env:"SCHEDULE_WORKER_INTERVAL,required"`
	ScheduleWorkerInterval    time.Duration

	// Interval for schedule poller
	UpcomingEventsWorkerIntervalStr string `env:"UPCOMING_EVENTS_WORKER_INTERVAL,required"`
	UpcomingEventsWorkerInterval    time.Duration

	// How many upcoming events to show per group in the pinned message
	UpcomingEventsShowCount int `env:"UPCOMING_EVENTS_SHOW_COUNT,default=5"`

	// Optional Gemini AI config for smart update summaries
	GeminiAPIKey string `env:"GEMINI_API_KEY"`
	GeminiModel  string `env:"GEMINI_MODEL"`
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

	duration, err := time.ParseDuration(cfg.ScheduleWorkerIntervalStr)
	if err != nil {
		return cfg, fmt.Errorf("invalid duration: %w", err)
	}
	cfg.ScheduleWorkerInterval = duration

	duration, err = time.ParseDuration(cfg.UpcomingEventsWorkerIntervalStr)
	if err != nil {
		return cfg, fmt.Errorf("invalid duration: %w", err)
	}
	cfg.UpcomingEventsWorkerInterval = duration

	return cfg, err
}
