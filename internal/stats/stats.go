// Package stats provides log entry statistics and summary computation.
package stats

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Summary holds aggregated statistics over a set of log entries.
type Summary struct {
	Total     int
	ByLevel   map[string]int
	FirstTime *time.Time
	LastTime  *time.Time
	Fields    map[string]map[string]int // field -> value -> count
}

// Compute calculates statistics from a slice of log entries.
func Compute(entries []parser.Entry) Summary {
	s := Summary{
		ByLevel: make(map[string]int),
		Fields:  make(map[string]map[string]int),
	}

	for _, e := range entries {
		s.Total++

		level := string(e.Level)
		s.ByLevel[level]++

		if !e.Timestamp.IsZero() {
			if s.FirstTime == nil || e.Timestamp.Before(*s.FirstTime) {
				t := e.Timestamp
				s.FirstTime = &t
			}
			if s.LastTime == nil || e.Timestamp.After(*s.LastTime) {
				t := e.Timestamp
				s.LastTime = &t
			}
		}

		for k, v := range e.Fields {
			if s.Fields[k] == nil {
				s.Fields[k] = make(map[string]int)
			}
			s.Fields[k][v]++
		}
	}

	return s
}

// Duration returns the time span between first and last entry, or zero.
func (s Summary) Duration() time.Duration {
	if s.FirstTime == nil || s.LastTime == nil {
		return 0
	}
	return s.LastTime.Sub(*s.FirstTime)
}
