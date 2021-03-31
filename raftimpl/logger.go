package raftimpl

import (
	"github.com/sirupsen/logrus"
	"os"
)

type SimpleLogger struct {
	log *logrus.Logger
}

func NewLogger() *SimpleLogger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.TraceLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
		PadLevelText: true,
	})
	return &SimpleLogger{log: log}
}

func (l *SimpleLogger) Trace(msg string) {
	l.log.Trace(msg)
}

func (l *SimpleLogger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l *SimpleLogger) Info(msg string) {
	l.log.Info(msg)
}

func (l *SimpleLogger) Warn(msg string) {
	l.log.Warn(msg)
}

func (l *SimpleLogger) Error(msg string) {
	l.log.Error(msg)
}

