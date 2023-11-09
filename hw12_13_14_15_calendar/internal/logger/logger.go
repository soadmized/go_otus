package logger

import (
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	l zerolog.Logger
}

func New(level string) *Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		panic(err)
	}

	logger := zerolog.New(os.Stdout).Level(lvl)

	return &Logger{logger}
}

func (l Logger) Info(msg string) {
	l.l.Info().Msg(msg)
}

func (l Logger) Error(msg string) {
	l.l.Error().Msg(msg)
}
