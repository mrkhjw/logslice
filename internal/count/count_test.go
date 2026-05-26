package count_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/count"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(level, msg string, fields map[string]string) parser.Entry {
	if fields == nil {
		fields = map[string]string{}
	}
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_EmptyInput(t *testing.T) {
	results := count.Apply(nil, count.DefaultOptions())
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestApply_GroupByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", nil),
		makeEntry("error", "c", nil),
	}
	results := count.Apply(entries, count.DefaultOptions())
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
		makeEntry("info", "a", map[string]string{"service": "api"}),
		makeEntry("info", "b", map[string]string{"service": "api"}),
		makeEntry("info", "c", map[string]string{"service": "worker"}),
	}
	opts := count.Options{Field: "service"}
	results := count.Apply(entries, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
	if results[0].Value != "api" || results[0].Count != 2 {
		t.Errorf("expected api=2, got %s=%d", results[0].Value, results[0].Count)
	}
}

func TestApply_TopN_Limits(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", nil),
		makeEntry("warn", "c", nil),
		makeEntry("error", "d", nil),
	}
	opts := count.Options{Field: "level", TopN: 2}
	results := count.Apply(entries, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results with TopN=2, got %d", len(results))
	}
}

func TestApply_MissingField_EmptyKey(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", map[string]string{"svc": "x"}),
	}
	opts := count.Options{Field: "svc"}
	results := count.Apply(entries, opts)
	// expect "x"=1 and ""=1
	if len(results) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(results))
	}
}
