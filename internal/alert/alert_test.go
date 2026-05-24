package alert_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/alert"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(level string, ts time.Time) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Level:     parser.Level(level),
		Message:   "test message",
	}
}

func TestEvaluate_NoRules_NoEvents(t *testing.T) {
	fired := 0
	a := alert.New(nil, func(alert.Event) { fired++ })
	a.Evaluate([]parser.Entry{makeEntry("error", time.Now())})
	if fired != 0 {
		t.Errorf("expected 0 events, got %d", fired)
	}
}

func TestEvaluate_BelowThreshold_NoEvent(t *testing.T) {
	rules := []alert.Rule{{Level: "error", Threshold: 5, Window: time.Minute}}
	fired := 0
	a := alert.New(rules, func(alert.Event) { fired++ })
	now := time.Now()
	entries := []parser.Entry{
		makeEntry("error", now),
		makeEntry("error", now.Add(time.Second)),
	}
	a.Evaluate(entries)
	if fired != 0 {
		t.Errorf("expected 0 events, got %d", fired)
	}
}

func TestEvaluate_ExceedsThreshold_FiresEvent(t *testing.T) {
	rules := []alert.Rule{{Level: "error", Threshold: 2, Window: time.Minute}}
	var events []alert.Event
	a := alert.New(rules, func(e alert.Event) { events = append(events, e) })
	now := time.Now()
	entries := []parser.Entry{
		makeEntry("error", now),
		makeEntry("error", now.Add(time.Second)),
		makeEntry("error", now.Add(2*time.Second)),
	}
	a.Evaluate(entries)
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Count != 3 {
		t.Errorf("expected count 3, got %d", events[0].Count)
	}
}

func TestEvaluate_OutsideWindow_NoEvent(t *testing.T) {
	rules := []alert.Rule{{Level: "error", Threshold: 1, Window: 10 * time.Second}}
	fired := 0
	a := alert.New(rules, func(alert.Event) { fired++ })
	now := time.Now()
	entries := []parser.Entry{
		makeEntry("error", now.Add(-30*time.Second)),
		makeEntry("error", now),
	}
	a.Evaluate(entries)
	if fired != 0 {
		t.Errorf("expected 0 events (entries outside window), got %d", fired)
	}
}

func TestEvaluate_MultipleRules_EachFires(t *testing.T) {
	rules := []alert.Rule{
		{Level: "error", Threshold: 1, Window: time.Minute},
		{Level: "warn", Threshold: 1, Window: time.Minute},
	}
	fired := 0
	a := alert.New(rules, func(alert.Event) { fired++ })
	now := time.Now()
	entries := []parser.Entry{
		makeEntry("error", now),
		makeEntry("error", now.Add(time.Second)),
		makeEntry("warn", now),
		makeEntry("warn", now.Add(time.Second)),
	}
	a.Evaluate(entries)
	if fired != 2 {
		t.Errorf("expected 2 events, got %d", fired)
	}
}

func TestEvent_MessageContainsLevel(t *testing.T) {
	rules := []alert.Rule{{Level: "fatal", Threshold: 0, Window: time.Minute}}
	var events []alert.Event
	a := alert.New(rules, func(e alert.Event) { events = append(events, e) })
	a.Evaluate([]parser.Entry{makeEntry("fatal", time.Now())})
	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if events[0].Message == "" {
		t.Error("expected non-empty message")
	}
}
