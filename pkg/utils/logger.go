package utils

import (
	"log"
	"os"
	"path/filepath"
)

func NewLogger() (*os.File, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execPath)

	logFilePath := filepath.Join(execDir, "logs", "bsync_cli_main.log")
	logDir := filepath.Dir(logFilePath)

	err = os.MkdirAll(logDir, 0777)
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	log.SetOutput(logFile)

	return logFile, nil
}
