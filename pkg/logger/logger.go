package logger

import (
	"time"

	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339

	log.Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        colorable.NewColorableStdout(),
		TimeFormat: "15:04:05",
		NoColor:    false,
	}).With().Timestamp().Logger()
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Error(msg string, err error) {
	log.Error().Err(err).Msg(msg)
}

func Debug(msg string) {
	log.Debug().Msg(msg)
}
