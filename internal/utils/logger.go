package utils

import (
	"io"
	"log"
	"os"
	"unique"
)

// Log levels
const (
	DebugLevel = iota
	InfoLevel
	SuccessLevel
	WarningLevel
	ErrorLevel
)

type Logger struct {
	Level       int
	lastLogged  unique.Handle[string]
	debugLogger *log.Logger
	infoLogger  *log.Logger
	succeLogger *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

var logger *Logger
var logFile *os.File

func init() {
	logFile, err := os.OpenFile("aki_ddos.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)

	logger = &Logger{
		Level:       InfoLevel,
		debugLogger: log.New(multiWriter, Gray+"[=] DEBUG: "+White, log.LstdFlags),
		infoLogger:  log.New(multiWriter, Cyan+"[*] INFO: "+White, log.LstdFlags),
		succeLogger: log.New(multiWriter, Green+"[+] SUCCESS: "+White, log.LstdFlags),
		warnLogger:  log.New(multiWriter, Yellow+"[/] WARN: "+White, log.LstdFlags),
		errorLogger: log.New(multiWriter, Red+"[//] ERROR: "+White, log.LstdFlags),
	}
}

func CloseLogFile() {
	if logFile != nil {
		logFile.Close()
	}
}

// Set log level
func SetLevel(level int) {
	logger.Level = level
}

func Debug(message string) {
	if logger.Level <= DebugLevel {
		logger.debugLogger.Println(message)
	}
}

func Info(message string) {
	if logger.Level <= InfoLevel {
		logger.infoLogger.Println(message)
	}
}

func Success(message string) {
	if logger.Level <= SuccessLevel {
		logger.succeLogger.Println(message)
	}
}

func Warning(message string) {
	if logger.Level <= WarningLevel {
		logger.warnLogger.Println(message)
	}
}

func Error(message string) {
	if logger.Level <= ErrorLevel {
		logger.errorLogger.Println(message)
	}
}
