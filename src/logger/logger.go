package logger

import (
	"log"
	"os"
	"path/filepath"
)

var (
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger
)

// creates directory for logs and sets output to file
func createLogsfile() *os.File {
	logsDir := filepath.Join(".", "logs")
	err := os.MkdirAll(logsDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logfile, err := os.Create(filepath.Join(logsDir, "logs.log"))
	log.SetOutput(logfile)

	return logfile
}

// creates new custom loggers
func setUpLoggers(logfile *os.File) {
	infoLogger = log.New(logfile, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(logfile, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(logfile, "ERROR: ", log.Ldate|log.Ltime)
}

func init() {
	logfile := createLogsfile()
	setUpLoggers(logfile)
}

func LogInfo(message ...interface{}) {
	infoLogger.Println(message...)
}

func LogWarning(message ...interface{}) {
	warningLogger.Println(message...)
}

func LogError(isFatal bool, message ...interface{}) {
	if isFatal {
		errorLogger.Fatal(message...)
	}
	errorLogger.Println(message...)
}
