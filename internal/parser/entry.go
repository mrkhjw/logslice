package parser

import (
	"fmt"
	"time"
)

// Level represents a log severity level.
type Level string

const (
	LevelDebug Level = "DEBUG"
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelError Level = "ERROR"
	LevelFatal Level = "FATAL"
	LevelUnknown Level = "UNKNOWN"
)

// Entry represents a single parsed log line.
type Entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     Level             `json:"level"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	Raw       string            `json:"raw"`
}

// String returns a human-readable representation of the log entry.
func (e *Entry) String() string {
	return fmt.Sprintf("[%s] %s %s",
		e.Timestamp.Format(time.RFC3339),
		e.Level,
		e.Message,
	)
}

// ParseLevel converts a string to a Level, returning LevelUnknown if unrecognized.
func ParseLevel(s string) Level {
	switch Level(s) {
	case LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal:
		return Level(s)
	default:
		return LevelUnknown
	}
}
