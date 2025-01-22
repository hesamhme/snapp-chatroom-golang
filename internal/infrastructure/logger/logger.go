package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Init initializes the global logger configuration
func Init() {
	logFile, err := os.OpenFile("log/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	logrus.SetOutput(logFile)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Info("Logger initialized successfully")
}
