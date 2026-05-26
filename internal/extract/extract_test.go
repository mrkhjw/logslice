package extract_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/extract"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(level, msg string, fields map[string]string) parser.Entry {
	e := parser.Entry{
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]string),
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

func TestApply_EmptyInput(t *testing.T) {
	result := extract.Apply(nil, extract.DefaultOptions())
	if len(result) != 0 {
		t.Fatalf("expected empty, got %d", len(result))
	}
}

func TestApply_IncludesMessageByDefault(t *testing.T) {
	entries := []parser.Entry{makeEntry("info", "hello world", nil)}
	result := extract.Apply(entries, extract.DefaultOptions())
	if result[0].Fields["message"] != "hello world" {
		t.Errorf("expected message field, got %q", result[0].Fields["message"])
	}
}

func TestApply_IncludesLevelByDefault(t *testing.T) {
	entries := []parser.Entry{makeEntry("error", "boom", nil)}
	result := extract.Apply(entries, extract.DefaultOptions())
	if result[0].Fields["level"] != "error" {
		t.Errorf("expected level=error, got %q", result[0].Fields["level"])
	}
}

func TestApply_ExtractsNamedField(t *testing.T) {
	entries := []parser.Entry{makeEntry("info", "msg", map[string]string{"user": "alice"})}
	opts := extract.DefaultOptions()
	opts.Fields = map[string]string{"user": ""}
	result := extract.Apply(entries, opts)
	if result[0].Fields["user"] != "alice" {
		t.Errorf("expected user=alice, got %q", result[0].Fields["user"])
	}
}

func TestApply_RenamesField(t *testing.T) {
	entries := []parser.Entry{makeEntry("info", "msg", map[string]string{"usr": "bob"})}
	opts := extract.DefaultOptions()
	opts.Fields = map[string]string{"usr": "user"}
	result := extract.Apply(entries, opts)
	if result[0].Fields["user"] != "bob" {
		t.Errorf("expected user=bob, got %q", result[0].Fields["user"])
	}
	if _, ok := result[0].Fields["usr"]; ok {
		t.Error("original key should not be present")
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	entries := []parser.Entry{makeEntry("info", "msg", nil)}
	opts := extract.DefaultOptions()
	opts.Fields = map[string]string{"nonexistent": "out"}
	result := extract.Apply(entries, opts)
	if _, ok := result[0].Fields["out"]; ok {
		t.Error("missing source field should not produce output key")
	}
}

func TestApply_DisableBuiltins(t *testing.T) {
	entries := []parser.Entry{makeEntry("warn", "msg", nil)}
	opts := extract.Options{
		Fields:           map[string]string{},
		IncludeMessage:   false,
		IncludeLevel:     false,
		IncludeTimestamp: false,
	}
	result := extract.Apply(entries, opts)
	if len(result[0].Fields) != 0 {
		t.Errorf("expected empty fields, got %v", result[0].Fields)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	origFields := map[string]string{"k": "v"}
	entries := []parser.Entry{makeEntry("info", "msg", origFields)}
	opts := extract.DefaultOptions()
	opts.Fields = map[string]string{"k": "key"}
	extract.Apply(entries, opts)
	if _, ok := entries[0].Fields["key"]; ok {
		t.Error("original entry should not be mutated")
	}
}
