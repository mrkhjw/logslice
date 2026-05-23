package redact_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/redact"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_RedactsField(t *testing.T) {
	r := redact.New([]redact.Rule{{Field: "password"}})
	e := makeEntry("user login", map[string]string{"user": "alice", "password": "s3cr3t"})
	out := r.Apply(e)

	if out.Fields["password"] != "***REDACTED***" {
		t.Errorf("expected password to be redacted, got %q", out.Fields["password"])
	}
	if out.Fields["user"] != "alice" {
		t.Errorf("user field should be unchanged, got %q", out.Fields["user"])
	}
}

func TestApply_CustomMask(t *testing.T) {
	r := redact.New([]redact.Rule{{Field: "token", Mask: "[HIDDEN]"}})
	e := makeEntry("auth", map[string]string{"token": "abc123"})
	out := r.Apply(e)

	if out.Fields["token"] != "[HIDDEN]" {
		t.Errorf("expected [HIDDEN], got %q", out.Fields["token"])
	}
}

func TestApply_PatternInMessage(t *testing.T) {
	rule := redact.MustCompile(`\b\d{4}-\d{4}-\d{4}-\d{4}\b`, "****-****-****-****")
	r := redact.New([]redact.Rule{rule})
	e := makeEntry("card 1234-5678-9012-3456 charged", nil)
	out := r.Apply(e)

	expected := "card ****-****-****-**** charged"
	if out.Message != expected {
		t.Errorf("expected %q, got %q", expected, out.Message)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	r := redact.New([]redact.Rule{{Field: "secret"}})
	e := makeEntry("test", map[string]string{"secret": "value"})
	_ = r.Apply(e)

	if e.Fields["secret"] != "value" {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_MissingField_NoError(t *testing.T) {
	r := redact.New([]redact.Rule{{Field: "nonexistent"}})
	e := makeEntry("hello", map[string]string{"key": "val"})
	out := r.Apply(e)

	if out.Fields["key"] != "val" {
		t.Error("unrelated fields should be preserved")
	}
}

func TestApplyAll_ProcessesSlice(t *testing.T) {
	r := redact.New([]redact.Rule{{Field: "pw"}})
	entries := []parser.Entry{
		makeEntry("a", map[string]string{"pw": "x"}),
		makeEntry("b", map[string]string{"pw": "y"}),
	}
	out := r.ApplyAll(entries)

	if len(out) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(out))
	}
	for i, e := range out {
		if e.Fields["pw"] != "***REDACTED***" {
			t.Errorf("entry %d: expected redacted pw, got %q", i, e.Fields["pw"])
		}
	}
}
