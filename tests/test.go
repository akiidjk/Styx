package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type Colors struct {
	Reset   string
	Red     string
	Green   string
	Yellow  string
	Blue    string
	Magenta string
	Cyan    string
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
	White:   "\033[37m",
}

// Possible levels:
// panic (zerolog.PanicLevel, 5)
// fatal (zerolog.FatalLevel, 4)
// error (zerolog.ErrorLevel, 3)
// warn (zerolog.WarnLevel, 2)
// info (zerolog.InfoLevel, 1)
// debug (zerolog.DebugLevel, 0)
// trace (zerolog.TraceLevel, -1)

var logger zerolog.Logger
var logFilename string = time.Now().Local().Format("2006-01-02_15-04-05") + ".log"
var logFile *os.File

const logDir string = "styx"

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

	// Some customizations
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger = zerolog.New(io.Writer(logFile)).With().Caller().Timestamp().Logger()

}

func SetLevel(level zerolog.Level) {
	zerolog.SetGlobalLevel(level)
}
