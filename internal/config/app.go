package config

import (
	"bytes"
	"os"

	"github.com/DeanPDX/dotconfig"
)

type AppConfig struct {
	// ID of google spreadsheets document with schedule
	SpreadsheetsDocumentID string `env:"SCHEDULE_DOCUMENT_ID,required"`
}

func LoadAppConfig(src LoadSource) (AppConfig, error) {
	switch src {
	case EnviromentVariablesSource:
		var buf bytes.Buffer

		for _, e := range os.Environ() {
			// e = "KEY=VALUE"
			buf.WriteString(e)
			buf.WriteByte('\n')
		}

		return dotconfig.FromReader[AppConfig](&buf)
	case EnvFileSource:
		return dotconfig.FromFileName[AppConfig](".env")
	default:
		panic("wrong load source option")
	}
}
