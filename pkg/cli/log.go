package cli

import (
	"github.com/rs/zerolog"
)

func ApplicationSuccefulStopped(l zerolog.Logger) {
	l.Info().Msg("Application succeful stopped")
}

func ApplicationStartFailed(l zerolog.Logger, err error) {
	l.Fatal().Err(err).Msg("Application start failed")
}

func ApplicationStarted(l zerolog.Logger) {
	l.Info().Msg("Application succeful started")
}
