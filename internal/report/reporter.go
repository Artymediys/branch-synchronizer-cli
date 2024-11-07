package report

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func NewReporter() (*os.File, error) {
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execDir := filepath.Dir(execPath)

	reportDate := time.Now().Format(time.DateOnly)
	fileName := fmt.Sprintf("bsync_cli_result_%s.log", reportDate)

	reportFilePath := filepath.Join(execDir, "logs", fileName)
	logDir := filepath.Dir(reportFilePath)

	err = os.MkdirAll(logDir, 0777)
	if err != nil {
		return nil, err
	}

	reportFile, err := os.OpenFile(reportFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return reportFile, nil
}
