package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() error {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	logFile, err := os.OpenFile("logs/app.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	InfoLogger = log.New(logFile, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(logFile, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}
