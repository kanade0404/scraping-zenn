package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Infof(format string, msgs ...interface{}) {
	log.Info().Msgf(format, msgs)
}

func Fatal(msg string) {
	log.Fatal().Msg(msg)
	os.Exit(1)
}

func Fatalf(format string, msgs ...interface{}) {
	log.Fatal().Msgf(format, msgs)
	os.Exit(1)
}

func Error(msg string) {
	log.Error().Msg(msg)
}

func Errorf(format string, msgs ...interface{}) {
	log.Error().Msgf(format, msgs)
}
