package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
	"unique"
)

type Colors struct {
	Reset   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
	Gray    string
	White   string
}

var DefaultColors = Colors{
	Reset:   "\033[0m",
	Red:     "\033[31m",
	Green:   "\033[32m",
	Yellow:  "\033[33m",
	Blue:    "\033[34m",
	Magenta: "\033[35m",
	Cyan:    "\033[36m",
	Gray:    "\033[37m",
	White:   "\033[97m",
}

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
	fatalLogger *log.Logger
}

var logDir string = "styx"
var logFilename string = time.Now().Local().Format("02-01-2006_15-04-05") + ".log"
var logger *Logger
var logFile *os.File

func init() {

	eiud := os.Geteuid()
	if eiud != 0 {
		log.Fatal("Remember to run the program with `sudo` or with root")
		os.Exit(1)
	}

	newpath := filepath.Join("/", "var", "log", logDir)
	err := os.MkdirAll(newpath, os.ModePerm)
	pathLogs := filepath.Join(newpath, logFilename)

	logFile, err := os.OpenFile(pathLogs, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)

	logger = &Logger{
		Level:       InfoLevel,
		debugLogger: log.New(multiWriter, DefaultColors.Gray+"[=] DEBUG: "+DefaultColors.White, log.LstdFlags),
		infoLogger:  log.New(multiWriter, DefaultColors.Cyan+"[*] INFO: "+DefaultColors.White, log.LstdFlags),
		succeLogger: log.New(multiWriter, DefaultColors.Green+"[+] SUCCESS: "+DefaultColors.White, log.LstdFlags),
		warnLogger:  log.New(multiWriter, DefaultColors.Yellow+"[/] WARN: "+DefaultColors.White, log.LstdFlags),
		fatalLogger: log.New(multiWriter, DefaultColors.Red+"[//] ERROR: "+DefaultColors.White, log.LstdFlags),
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

func Fatal(message ...any) {
	if logger.Level <= ErrorLevel {
		logger.fatalLogger.Println(fmt.Sprint(message...))
	}
}
