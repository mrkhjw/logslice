package mask_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/mask"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_NoOptions_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello", map[string]string{"user": "alice"})}
	out := mask.Apply(entries, mask.Options{})
	if out[0].Fields["user"] != "alice" {
		t.Errorf("expected alice, got %s", out[0].Fields["user"])
	}
}

func TestApply_MasksNamedField(t *testing.T) {
	entries := []parser.Entry{makeEntry("login", map[string]string{"password": "secret", "user": "bob"})}
	out := mask.Apply(entries, mask.Options{Fields: []string{"password"}})
	if out[0].Fields["password"] != mask.DefaultMask {
		t.Errorf("expected masked password, got %s", out[0].Fields["password"])
	}
	if out[0].Fields["user"] != "bob" {
		t.Errorf("unrelated field should be unchanged")
	}
}

func TestApply_CustomMask(t *testing.T) {
	entries := []parser.Entry{makeEntry("x", map[string]string{"token": "abc123"})}
	out := mask.Apply(entries, mask.Options{Fields: []string{"token"}, Mask: "[REDACTED]"})
	if out[0].Fields["token"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %s", out[0].Fields["token"])
	}
}

func TestApply_MaskMessage(t *testing.T) {
	entries := []parser.Entry{makeEntry("sensitive data here", nil)}
	out := mask.Apply(entries, mask.Options{MaskMessage: true})
	if out[0].Message != mask.DefaultMask {
		t.Errorf("expected masked message, got %s", out[0].Message)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := makeEntry("msg", map[string]string{"secret": "value"})
	entries := []parser.Entry{original}
	mask.Apply(entries, mask.Options{Fields: []string{"secret"}})
	if entries[0].Fields["secret"] != "value" {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_CaseInsensitiveFieldName(t *testing.T) {
	entries := []parser.Entry{makeEntry("x", map[string]string{"APIKey": "key123"})}
	out := mask.Apply(entries, mask.Options{Fields: []string{"apikey"}})
	if out[0].Fields["APIKey"] != mask.DefaultMask {
		t.Errorf("field matching should be case-insensitive, got %s", out[0].Fields["APIKey"])
	}
}

func TestApply_EmptyEntries(t *testing.T) {
	out := mask.Apply(nil, mask.Options{Fields: []string{"x"}})
	if len(out) != 0 {
		t.Errorf("expected empty output for nil input")
	}
}
