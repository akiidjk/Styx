package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"unique"

	"github.com/akiidjk/fw-ngfw/internal/utils"
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
	logFile, err := os.OpenFile("logs/styx.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)

	logger = &Logger{
		Level:       InfoLevel,
		debugLogger: log.New(multiWriter, utils.Gray+"[=] DEBUG: "+utils.White, log.LstdFlags),
		infoLogger:  log.New(multiWriter, utils.Cyan+"[*] INFO: "+utils.White, log.LstdFlags),
		succeLogger: log.New(multiWriter, utils.Green+"[+] SUCCESS: "+utils.White, log.LstdFlags),
		warnLogger:  log.New(multiWriter, utils.Yellow+"[/] WARN: "+utils.White, log.LstdFlags),
		errorLogger: log.New(multiWriter, utils.Red+"[//] ERROR: "+utils.White, log.LstdFlags),
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

func Debug(message ...any) {
	if logger.Level <= DebugLevel {
		logger.debugLogger.Println(fmt.Sprint(message...))
	}
}

func Info(message ...any) {
	if logger.Level <= InfoLevel {
		logger.infoLogger.Println(fmt.Sprint(message...))
	}
}

func Success(message ...any) {
	if logger.Level <= SuccessLevel {
		logger.succeLogger.Println(fmt.Sprint(message...))
	}
}

func Warning(message ...any) {
	if logger.Level <= WarningLevel {
		logger.warnLogger.Println(fmt.Sprint(message...))
	}
}

func Error(message ...any) {
	if logger.Level <= ErrorLevel {
		logger.errorLogger.Println(fmt.Sprint(message...))
	}
}
