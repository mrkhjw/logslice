package top_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/top"
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
	results := top.Apply(nil, top.DefaultOptions())
	if len(results) != 0 {
		t.Fatalf("expected empty results, got %d", len(results))
	}
}

func TestApply_TopByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("error", "a", nil),
		makeEntry("error", "b", nil),
		makeEntry("warn", "c", nil),
		makeEntry("info", "d", nil),
		makeEntry("error", "e", nil),
	}
	opts := top.Options{Field: "level", N: 2}
	results := top.Apply(entries, opts)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Value != "error" || results[0].Count != 3 {
		t.Errorf("expected error×3, got %s×%d", results[0].Value, results[0].Count)
	}
	if results[1].Value != "warn" && results[1].Value != "info" {
		t.Errorf("unexpected second result: %s", results[1].Value)
	}
}

func TestApply_TopByCustomField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "x", map[string]string{"host": "a"}),
		makeEntry("info", "x", map[string]string{"host": "b"}),
		makeEntry("info", "x", map[string]string{"host": "a"}),
	}
	opts := top.Options{Field: "host", N: 10}
	results := top.Apply(entries, opts)
	if results[0].Value != "a" || results[0].Count != 2 {
		t.Errorf("expected a×2, got %s×%d", results[0].Value, results[0].Count)
	}
}

func TestApply_ZeroN_ReturnsAll(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "x", nil),
		makeEntry("warn", "y", nil),
		makeEntry("error", "z", nil),
	}
	opts := top.Options{Field: "level", N: 0}
	results := top.Apply(entries, opts)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestApply_MissingField_CountedAsEmpty(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("info", "a", nil),
		makeEntry("info", "b", nil),
	}
	opts := top.Options{Field: "nonexistent", N: 5}
	results := top.Apply(entries, opts)
	if len(results) != 1 || results[0].Value != "" || results[0].Count != 2 {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := top.DefaultOptions()
	if opts.Field != "level" {
		t.Errorf("expected field 'level', got %q", opts.Field)
	}
	if opts.N != 10 {
		t.Errorf("expected N=10, got %d", opts.N)
	}
}
