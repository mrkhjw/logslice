package aggregate_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/aggregate"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(level, message string, fields map[string]any) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Fields:    fields,
	}
}

func TestApply_GroupByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", nil),
		makeEntry("error", "c", nil),
	}
	results := aggregate.Apply(entries, aggregate.Options{Field: "level"})
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	if results[0].Value != "info" || results[0].Count != 2 {
		t.Errorf("expected info=2, got %s=%d", results[0].Value, results[0].Count)
	}
	if results[1].Value != "error" || results[1].Count != 1 {
		t.Errorf("expected error=1, got %s=%d", results[1].Value, results[1].Count)
	}
}

func TestApply_GroupByCustomField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "x", map[string]any{"host": "web-1"}),
		makeEntry("info", "y", map[string]any{"host": "web-1"}),
		makeEntry("info", "z", map[string]any{"host": "db-1"}),
	}
	results := aggregate.Apply(entries, aggregate.Options{Field: "host"})
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	if results[0].Value != "web-1" || results[0].Count != 2 {
		t.Errorf("unexpected top result: %+v", results[0])
	}
}

func TestApply_TopN(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("warn", "b", nil),
		makeEntry("error", "c", nil),
		makeEntry("debug", "d", nil),
	}
	results := aggregate.Apply(entries, aggregate.Options{Field: "level", TopN: 2})
	if len(results) != 2 {
		t.Fatalf("expected 2 results with TopN=2, got %d", len(results))
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	results := aggregate.Apply(nil, aggregate.Options{Field: "level"})
	if len(results) != 0 {
		t.Errorf("expected empty results, got %d", len(results))
	}
}

func TestApply_MissingField_GroupedAsEmpty(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", nil),
	}
	results := aggregate.Apply(entries, aggregate.Options{Field: "nonexistent"})
	if len(results) != 1 {
		t.Fatalf("expected 1 group for missing field, got %d", len(results))
	}
	if results[0].Value != "" {
		t.Errorf("expected empty string value, got %q", results[0].Value)
	}
	if results[0].Count != 2 {
		t.Errorf("expected count 2, got %d", results[0].Count)
	}
}
