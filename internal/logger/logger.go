package logger

import (
	"os"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

func SetupLogger() {
	Log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}
