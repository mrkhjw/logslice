package filter

import (
	"time"

	"github.com/user/logslice/internal/parser"
)

// Options holds the filtering criteria for log entries.
type Options struct {
	Level     string
	Since     time.Time
	Until     time.Time
	FieldKey  string
	FieldVal  string
}

// Filter applies the given options to a slice of log entries
// and returns only those that match all specified criteria.
func Filter(entries []parser.Entry, opts Options) []parser.Entry {
	var result []parser.Entry
	for _, e := range entries {
		if !matchLevel(e, opts.Level) {
			continue
		}
		if !matchTimeRange(e, opts.Since, opts.Until) {
			continue
		}
		if !matchField(e, opts.FieldKey, opts.FieldVal) {
			continue
		}
		result = append(result, e)
	}
	return result
}

func matchLevel(e parser.Entry, level string) bool {
	if level == "" {
		return true
	}
	return string(e.Level) == level
}

func matchTimeRange(e parser.Entry, since, until time.Time) bool {
	if !since.IsZero() && e.Timestamp.Before(since) {
		return false
	}
	if !until.IsZero() && e.Timestamp.After(until) {
		return false
	}
	return true
}

func matchField(e parser.Entry, key, val string) bool {
	if key == "" {
		return true
	}
	v, ok := e.Fields[key]
	if !ok {
		return false
	}
	if val == "" {
		return true
	}
	return v == val
}
