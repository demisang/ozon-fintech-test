package logger

import "github.com/sirupsen/logrus"

func GetLogger() *logrus.Logger {
	return logrus.New()
}
