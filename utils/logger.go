package utils

import (
	"github.com/sirupsen/logrus"
	"os"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.Out = os.Stdout
	Logger.SetLevel(logrus.InfoLevel)

	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}