package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	WithFields(fields logrus.Fields) *logrus.Entry
}

type AppLogger struct {
	*logrus.Logger
}

func NewLogger() *AppLogger {
	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetFormatter(&logrus.JSONFormatter{})

	env := os.Getenv("APP_ENV")
	if env == "development" {
		l.SetLevel(logrus.DebugLevel)
	} else {
		l.SetLevel(logrus.InfoLevel)
	}

	return &AppLogger{l}
}
