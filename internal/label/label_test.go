package label_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/label"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg, level string, fields map[string]string) parser.Entry {
	if fields == nil {
		fields = make(map[string]string)
	}
	return parser.Entry{Timestamp: time.Now(), Level: level, Message: msg, Fields: fields}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello", "info", nil)}
	out := label.Apply(entries, label.DefaultOptions())
	if len(out) != 1 || out[0].Message != "hello" {
		t.Fatal("expected original entry")
	}
}

func TestApply_MatchMessage_AddsLabel(t *testing.T) {
	rule := label.Rule{Pattern: regexp.MustCompile(`timeout`), Label: "slow"}
	opts := label.DefaultOptions()
	opts.Rules = []label.Rule{rule}
	entries := []parser.Entry{makeEntry("connection timeout", "warn", nil)}
	out := label.Apply(entries, opts)
	if out[0].Fields["label"] != "slow" {
		t.Fatalf("expected label 'slow', got %q", out[0].Fields["label"])
	}
}

func TestApply_NoMatch_LabelAbsent(t *testing.T) {
	rule := label.Rule{Pattern: regexp.MustCompile(`panic`), Label: "critical"}
	opts := label.DefaultOptions()
	opts.Rules = []label.Rule{rule}
	entries := []parser.Entry{makeEntry("all good", "info", nil)}
	out := label.Apply(entries, opts)
	if _, ok := out[0].Fields["label"]; ok {
		t.Fatal("expected no label field")
	}
}

func TestApply_MatchNamedField(t *testing.T) {
	rule := label.Rule{Field: "service", Pattern: regexp.MustCompile(`^auth`), Label: "auth-svc"}
	opts := label.DefaultOptions()
	opts.Rules = []label.Rule{rule}
	entries := []parser.Entry{makeEntry("login", "info", map[string]string{"service": "auth-api"})}
	out := label.Apply(entries, opts)
	if out[0].Fields["label"] != "auth-svc" {
		t.Fatalf("expected 'auth-svc', got %q", out[0].Fields["label"])
	}
}

func TestApply_MultiMode_JoinsLabels(t *testing.T) {
	opts := label.DefaultOptions()
	opts.Multi = true
	opts.Rules = []label.Rule{
		{Pattern: regexp.MustCompile(`error`), Label: "err"},
		{Pattern: regexp.MustCompile(`disk`), Label: "storage"},
	}
	entries := []parser.Entry{makeEntry("disk error", "error", nil)}
	out := label.Apply(entries, opts)
	if out[0].Fields["label"] != "err,storage" {
		t.Fatalf("expected 'err,storage', got %q", out[0].Fields["label"])
	}
}

func TestApply_SingleMode_StopsAtFirstMatch(t *testing.T) {
	opts := label.DefaultOptions()
	opts.Multi = false
	opts.Rules = []label.Rule{
		{Pattern: regexp.MustCompile(`disk`), Label: "storage"},
		{Pattern: regexp.MustCompile(`error`), Label: "err"},
	}
	entries := []parser.Entry{makeEntry("disk error", "error", nil)}
	out := label.Apply(entries, opts)
	if out[0].Fields["label"] != "storage" {
		t.Fatalf("expected 'storage', got %q", out[0].Fields["label"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	rule := label.Rule{Pattern: regexp.MustCompile(`.*`), Label: "all"}
	opts := label.DefaultOptions()
	opts.Rules = []label.Rule{rule}
	original := makeEntry("msg", "info", nil)
	label.Apply([]parser.Entry{original}, opts)
	if _, ok := original.Fields["label"]; ok {
		t.Fatal("original entry was mutated")
	}
}
