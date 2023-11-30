package logutil

import "github.com/sirupsen/logrus"

var (
	Logger *logrus.Logger
)

func InitDefaultLogger(level logrus.Level) error {
	logger := logrus.New()

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})

	logger.AddHook(&DefaultHook{})

	Logger = logger
	return nil
}

type DefaultHook struct {
}

func (hook *DefaultHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *DefaultHook) Fire(entry *logrus.Entry) error {

	return nil
}
