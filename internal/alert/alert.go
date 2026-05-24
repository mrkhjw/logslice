// Package alert provides threshold-based alerting for log entry streams.
// It triggers callbacks when entry counts for a given level exceed a
// configured limit within a sliding time window.
package alert

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single alerting condition.
type Rule struct {
	// Level is the log level to monitor (e.g. "error", "fatal").
	Level string
	// Threshold is the maximum number of matching entries allowed in Window.
	Threshold int
	// Window is the duration of the sliding time window.
	Window time.Duration
}

// Event is emitted when a Rule's threshold is breached.
type Event struct {
	Rule    Rule
	Count   int
	Since   time.Time
	Message string
}

// Handler is called when an alert fires.
type Handler func(Event)

// Alerter evaluates log entries against a set of rules.
type Alerter struct {
	rules   []Rule
	handler Handler
	// buckets maps level -> slice of timestamps within the window
	buckets map[string][]time.Time
}

// New creates an Alerter with the given rules and handler.
func New(rules []Rule, handler Handler) *Alerter {
	return &Alerter{
		rules:   rules,
		handler: handler,
		buckets: make(map[string][]time.Time),
	}
}

// Evaluate processes a batch of entries and fires the handler for any
// rule whose threshold is exceeded within its window.
func (a *Alerter) Evaluate(entries []parser.Entry) {
	for _, e := range entries {
		lvl := string(e.Level)
		a.buckets[lvl] = append(a.buckets[lvl], e.Timestamp)
	}
	for _, rule := range a.rules {
		count, since := a.countInWindow(rule.Level, rule.Window)
		if count > rule.Threshold {
			a.handler(Event{
				Rule:    rule,
				Count:   count,
				Since:   since,
				Message: fmt.Sprintf("level %q exceeded threshold %d (got %d) in %s", rule.Level, rule.Threshold, count, rule.Window),
			})
		}
	}
}

// countInWindow returns the number of timestamps for level within the
// most recent window duration, and the start time of that window.
func (a *Alerter) countInWindow(level string, window time.Duration) (int, time.Time) {
	timestamps := a.buckets[level]
	if len(timestamps) == 0 {
		return 0, time.Time{}
	}
	cutoff := timestamps[len(timestamps)-1].Add(-window)
	count := 0
	for _, ts := range timestamps {
		if !ts.Before(cutoff) {
			count++
		}
	}
	return count, cutoff
}
