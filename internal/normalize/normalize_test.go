package normalize_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/normalize"
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

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{"host": "web1"})}
	out := normalize.Apply(entries, normalize.DefaultOptions())
	if len(out) != 1 || out[0].Fields["host"] != "web1" {
		t.Fatal("expected original entry unchanged")
	}
}

func TestApply_RenamesField(t *testing.T) {
	opts := normalize.Options{
		Rules: []normalize.Rule{{From: "hostname", To: "host"}},
	}
	entries := []parser.Entry{makeEntry(map[string]string{"hostname": "web1"})}
	out := normalize.Apply(entries, opts)
	if _, ok := out[0].Fields["hostname"]; ok {
		t.Error("old key 'hostname' should be removed")
	}
	if out[0].Fields["host"] != "web1" {
		t.Errorf("expected host=web1, got %q", out[0].Fields["host"])
	}
}

func TestApply_CaseInsensitiveMatch(t *testing.T) {
	opts := normalize.Options{
		Rules: []normalize.Rule{{From: "LEVEL", To: "severity"}},
	}
	entries := []parser.Entry{makeEntry(map[string]string{"level": "error"})}
	out := normalize.Apply(entries, opts)
	if out[0].Fields["severity"] != "error" {
		t.Errorf("expected severity=error, got %q", out[0].Fields["severity"])
	}
}

func TestApply_LowercaseValue(t *testing.T) {
	opts := normalize.Options{
		Rules: []normalize.Rule{{From: "env", To: "env", Lowercase: true}},
	}
	entries := []parser.Entry{makeEntry(map[string]string{"env": "PRODUCTION"})}
	out := normalize.Apply(entries, opts)
	if out[0].Fields["env"] != "production" {
		t.Errorf("expected env=production, got %q", out[0].Fields["env"])
	}
}

func TestApply_UppercaseValue(t *testing.T) {
	opts := normalize.Options{
		Rules: []normalize.Rule{{From: "status", To: "status", Uppercase: true}},
	}
	entries := []parser.Entry{makeEntry(map[string]string{"status": "ok"})}
	out := normalize.Apply(entries, opts)
	if out[0].Fields["status"] != "OK" {
		t.Errorf("expected status=OK, got %q", out[0].Fields["status"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	opts := normalize.Options{
		Rules: []normalize.Rule{{From: "src", To: "source"}},
	}
	origFields := map[string]string{"src": "app"}
	entries := []parser.Entry{makeEntry(origFields)}
	normalize.Apply(entries, opts)
	if _, ok := entries[0].Fields["src"]; !ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	opts := normalize.Options{
		Rules: []normalize.Rule{{From: "missing", To: "target"}},
	}
	entries := []parser.Entry{makeEntry(map[string]string{"other": "val"})}
	out := normalize.Apply(entries, opts)
	if _, ok := out[0].Fields["target"]; ok {
		t.Error("target field should not be added when source is missing")
	}
}
