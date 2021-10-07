package log

import "github.com/sirupsen/logrus"

func SetDefault(logger *logrus.Logger) *logrus.Logger {
	if logger == nil {
		return logrus.New()
	}
	return logger
}
