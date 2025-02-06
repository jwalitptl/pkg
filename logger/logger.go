package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	*zerolog.Logger
}

func NewLogger() *Logger {
	zerolog.TimeFieldFormat = time.RFC3339

	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}

	logger := zerolog.New(output).
		Level(zerolog.InfoLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	return &Logger{&logger}
}

func (l *Logger) Info(msg string, fields ...interface{}) {
	l.Logger.Info().Fields(fields).Msg(msg)
}

func (l *Logger) Error(err error, msg string, fields ...interface{}) {
	l.Logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Fatal(err error, msg string, fields ...interface{}) {
	l.Logger.Fatal().Err(err).Fields(fields).Msg(msg)
}

func (l *Logger) Debug(msg string, fields ...interface{}) {
	l.Logger.Debug().Fields(fields).Msg(msg)
}
