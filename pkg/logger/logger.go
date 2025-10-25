package logger

import (
	"log"
	"os"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
}

func New() *Logger {
	return &Logger{
		Info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) InfoLog(message string) {
	l.Info.Println(message)
}

func (l *Logger) ErrorLog(message string) {
	l.Error.Println(message)
}

func (l *Logger) DebugLog(message string) {
	l.Debug.Println(message)
}
