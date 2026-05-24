package flatten_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/flatten"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestApply_NoFields_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{})}
	result := flatten.Apply(entries, flatten.DefaultOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if len(result[0].Fields) != 0 {
		t.Errorf("expected empty fields, got %v", result[0].Fields)
	}
}

func TestApply_ScalarFields_Unchanged(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{
		"host":    "web-01",
		"service": "api",
	})}
	result := flatten.Apply(entries, flatten.DefaultOptions())
	if result[0].Fields["host"] != "web-01" {
		t.Errorf("host mismatch: %s", result[0].Fields["host"])
	}
	if result[0].Fields["service"] != "api" {
		t.Errorf("service mismatch: %s", result[0].Fields["service"])
	}
}

func TestApply_NestedValue_Flattened(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{
		"meta": "region=us-east zone=a",
	})}
	result := flatten.Apply(entries, flatten.DefaultOptions())
	fields := result[0].Fields
	if fields["meta.region"] != "us-east" {
		t.Errorf("expected meta.region=us-east, got %q", fields["meta.region"])
	}
	if fields["meta.zone"] != "a" {
		t.Errorf("expected meta.zone=a, got %q", fields["meta.zone"])
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Separator = "_"
	entries := []parser.Entry{makeEntry(map[string]string{
		"ctx": "user=alice role=admin",
	})}
	result := flatten.Apply(entries, opts)
	if result[0].Fields["ctx_user"] != "alice" {
		t.Errorf("expected ctx_user=alice, got %v", result[0].Fields)
	}
}

func TestApply_Prefix_Prepended(t *testing.T) {
	opts := flatten.DefaultOptions()
	opts.Prefix = "log"
	entries := []parser.Entry{makeEntry(map[string]string{
		"level": "warn",
	})}
	result := flatten.Apply(entries, opts)
	if result[0].Fields["log.level"] != "warn" {
		t.Errorf("expected log.level=warn, got %v", result[0].Fields)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := makeEntry(map[string]string{"k": "v"})
	entries := []parser.Entry{original}
	flatten.Apply(entries, flatten.DefaultOptions())
	if original.Fields["k"] != "v" {
		t.Error("original entry was mutated")
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	result := flatten.Apply(nil, flatten.DefaultOptions())
	if len(result) != 0 {
		t.Errorf("expected empty result, got %d entries", len(result))
	}
}
