// logger/logger.go
package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	InfoLogger  *logrus.Logger
	ErrorLogger *logrus.Logger
)

func Init() {
	// Create the info logger
	InfoLogger = logrus.New()
	infoFile, err := os.OpenFile("./apps/Accounts/logger/infolog.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		InfoLogger.Out = infoFile
	} else {
		InfoLogger.Info("Failed to log to file, using default stderr")
	}

	// Create the error logger
	ErrorLogger = logrus.New()
	errorFile, err := os.OpenFile("./apps/Accounts/logger/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		ErrorLogger.Out = errorFile
	} else {
		ErrorLogger.Error("Failed to log to file, using default stderr")
	}

	// Set log format to JSON for structured logging
	// InfoLogger.SetFormatter(&logrus.JSONFormatter{
	//  PrettyPrint: true,
	// })
	InfoLogger.SetFormatter(&logrus.JSONFormatter{})
	ErrorLogger.SetFormatter(&logrus.JSONFormatter{})

	InfoLogger.SetReportCaller(true)
	ErrorLogger.SetReportCaller(true)

	// Set log level
	InfoLogger.SetLevel(logrus.InfoLevel)
	ErrorLogger.SetLevel(logrus.ErrorLevel)
}
