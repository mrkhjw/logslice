package compact_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/compact"
	"github.com/logslice/logslice/internal/parser"
)

func makeEntry(fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestApply_NoOptions_DropEmptyByDefault(t *testing.T) {
	entry := makeEntry(map[string]string{"host": "srv1", "trace": ""})
	opts := compact.DefaultOptions()
	out := compact.Apply([]parser.Entry{entry}, opts)
	if _, ok := out[0].Fields["trace"]; ok {
		t.Error("expected empty field 'trace' to be dropped")
	}
	if out[0].Fields["host"] != "srv1" {
		t.Error("expected non-empty field 'host' to be retained")
	}
}

func TestApply_DropEmptyFalse_KeepsEmptyFields(t *testing.T) {
	entry := makeEntry(map[string]string{"host": "srv1", "trace": ""})
	opts := compact.DefaultOptions()
	opts.DropEmptyFields = false
	out := compact.Apply([]parser.Entry{entry}, opts)
	if _, ok := out[0].Fields["trace"]; !ok {
		t.Error("expected empty field 'trace' to be kept when DropEmptyFields=false")
	}
}

func TestApply_DropFields_RemovesNamed(t *testing.T) {
	entry := makeEntry(map[string]string{"host": "srv1", "secret": "abc123", "env": "prod"})
	opts := compact.DefaultOptions()
	opts.DropFields = []string{"secret"}
	out := compact.Apply([]parser.Entry{entry}, opts)
	if _, ok := out[0].Fields["secret"]; ok {
		t.Error("expected 'secret' to be dropped")
	}
	if out[0].Fields["host"] != "srv1" {
		t.Error("expected 'host' to be retained")
	}
}

func TestApply_KeepFields_ActsAsAllowlist(t *testing.T) {
	entry := makeEntry(map[string]string{"host": "srv1", "env": "prod", "extra": "noise"})
	opts := compact.DefaultOptions()
	opts.KeepFields = []string{"host", "env"}
	out := compact.Apply([]parser.Entry{entry}, opts)
	if _, ok := out[0].Fields["extra"]; ok {
		t.Error("expected 'extra' to be removed by keeplist")
	}
	if out[0].Fields["env"] != "prod" {
		t.Error("expected 'env' to be retained")
	}
}

func TestApply_PreservesCoreAttributes(t *testing.T) {
	entry := makeEntry(map[string]string{})
	out := compact.Apply([]parser.Entry{entry}, compact.DefaultOptions())
	if out[0].Level != "INFO" {
		t.Errorf("expected Level=INFO, got %s", out[0].Level)
	}
	if out[0].Message != "test message" {
		t.Errorf("expected Message preserved, got %s", out[0].Message)
	}
}

func TestApply_EmptyInput_ReturnsEmpty(t *testing.T) {
	out := compact.Apply([]parser.Entry{}, compact.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}
