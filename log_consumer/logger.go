package main

import (
	"io"
	"os"
	"time"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// KeyValueFormatter is a custom logrus formatter that formats log entries as key-value pairs in a single line
type KeyValueFormatter struct{}

// Format formats the log entry as key-value pairs in a single line
func (f *KeyValueFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format(time.RFC3339)
	return []byte(
		"from='" + entry.Data["from"].(string) + "'" +
			" log_level='" + entry.Data["log_level"].(string) + "'" +
			" msg='" + entry.Message + "'" +
			" time='" + timestamp + "'\n",
	), nil
}

// Logger is a struct that holds the logrus logger instance
type Logger struct {
	log *logrus.Logger
}

// NewLogger creates and returns a new Logger instance
func NewLogger(logFile string) *Logger {
	logger := &Logger{
		log: logrus.New(),
	}

	logger.init(logFile)

	return logger
}

// init initializes the logger with the specified log file
func (l *Logger) init(logFile string) {
	l.log.SetLevel(logrus.InfoLevel)
	l.log.SetFormatter(&KeyValueFormatter{})

	writer, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		hook := lfshook.NewHook(lfshook.WriterMap{
			logrus.InfoLevel:  writer,
			logrus.ErrorLevel: writer,
		}, &KeyValueFormatter{})
		l.log.AddHook(hook)

		// Disable console output
		// I don't want it to console it to the terminal
		l.log.SetOutput(io.Discard)
	} else {
		l.log.Error("Unable to open log file. Using default stderr")
	}
}

// LogEntry creates a log entry with the specified time, message, and from field
func (l *Logger) LogEntry(level string, message, from string, timestamp time.Time) {
	entry := logrus.NewEntry(l.log)
	entry.Time = timestamp
	entry.Message = message
	entry.Level = logrus.InfoLevel

	// Add custom fields
	entry.Data = make(logrus.Fields)
	entry.Data["from"] = from
	entry.Data["log_level"] = level

	entry.Logger = l.log

	l.log.WithFields(entry.Data).Log(entry.Level, entry.Message)
}

// Info logs informational messages with a specific time and from field
func (l *Logger) Info(message, from string, timestamp time.Time) {
	l.LogEntry("info", message, from, timestamp)
}

// Debug logs debug messages with a specific time and from field
func (l *Logger) Debug(message, from string, timestamp time.Time) {
	l.LogEntry("debug", message, from, timestamp)
}

// Error logs error messages with a specific time and from field
func (l *Logger) Error(message, from string, timestamp time.Time) {
	l.LogEntry("error", message, from, timestamp)
}
