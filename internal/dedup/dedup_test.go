package dedup_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/dedup"
	"github.com/yourorg/logslice/internal/parser"
)

var base = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func makeEntry(level, msg string, offset time.Duration) parser.Entry {
	return parser.Entry{
		Timestamp: base.Add(offset),
		Level:     level,
		Message:   msg,
		Fields:    map[string]string{},
	}
}

func TestDedup_NoDuplicates(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("INFO", "started", 0),
		makeEntry("INFO", "running", time.Second),
		makeEntry("WARN", "slow", 2*time.Second),
	}
	out := dedup.Dedup(entries, dedup.DefaultOptions())
	if len(out) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(out))
	}
}

func TestDedup_CollapsesConsecutiveDuplicates(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("ERROR", "connection refused", 0),
		makeEntry("ERROR", "connection refused", time.Second),
		makeEntry("ERROR", "connection refused", 2*time.Second),
		makeEntry("INFO", "recovered", 3*time.Second),
	}
	out := dedup.Dedup(entries, dedup.DefaultOptions())
	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	if out[0].Fields["dedup_count"] != "3" {
		t.Errorf("expected dedup_count=3, got %q", out[0].Fields["dedup_count"])
	}
}

func TestDedup_WindowExcludesFarEntries(t *testing.T) {
	opts := dedup.Options{Window: 30 * time.Second}
	entries := []parser.Entry{
		makeEntry("WARN", "timeout", 0),
		makeEntry("WARN", "timeout", 10*time.Second),
		makeEntry("WARN", "timeout", 60*time.Second), // outside window from prev
	}
	out := dedup.Dedup(entries, opts)
	if len(out) != 2 {
		t.Fatalf("expected 2 entries after window split, got %d", len(out))
	}
}

func TestDedup_MaxCountFlushes(t *testing.T) {
	opts := dedup.Options{Window: time.Hour, MaxCount: 2}
	entries := []parser.Entry{
		makeEntry("DEBUG", "poll", 0),
		makeEntry("DEBUG", "poll", time.Second),
		makeEntry("DEBUG", "poll", 2*time.Second),
	}
	out := dedup.Dedup(entries, opts)
	// MaxCount=2 means after 2 duplicates we flush, so we get 2 output entries
	if len(out) != 2 {
		t.Fatalf("expected 2 flushed entries, got %d", len(out))
	}
}

func TestDedup_EmptyInput(t *testing.T) {
	out := dedup.Dedup(nil, dedup.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output for nil input")
	}
}

func TestDedup_SingleEntry(t *testing.T) {
	entries := []parser.Entry{makeEntry("INFO", "hello", 0)}
	out := dedup.Dedup(entries, dedup.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if _, ok := out[0].Fields["dedup_count"]; ok {
		t.Error("single entry should not have dedup_count field")
	}
}
