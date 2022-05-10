package logger

import "github.com/sirupsen/logrus"

var LogrusLogger *logrus.Logger

func CreateLogger() {
	LogrusLogger = logrus.New()
}
