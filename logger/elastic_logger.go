package logger

import "github.com/sirupsen/logrus"

type ElasticLogger interface {
	Log(entity string, level string, message string)
}

type ElasticHook interface {
	Levels() []logrus.Level
	Fire(entry *logrus.Entry) error
}
