package classify_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/classify"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string, fields map[string]interface{}) parser.Entry {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello world", nil)}
	opts := classify.DefaultOptions()
	out := classify.Apply(entries, opts)
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if _, ok := out[0].Fields["category"]; ok {
		t.Error("expected no category field when no rules match and no default")
	}
}

func TestApply_MatchMessage_SetsCategory(t *testing.T) {
	pattern := regexp.MustCompile(`timeout`)
	opts := classify.DefaultOptions()
	opts.Rules = []classify.Rule{
		{Pattern: pattern, Category: "network"},
	}
	entries := []parser.Entry{
		makeEntry("connection timeout reached", nil),
		makeEntry("all good", nil),
	}
	out := classify.Apply(entries, opts)
	if got := out[0].Fields["category"]; got != "network" {
		t.Errorf("expected 'network', got %v", got)
	}
	if _, ok := out[1].Fields["category"]; ok {
		t.Error("second entry should not have category")
	}
}

func TestApply_MatchNamedField_SetsCategory(t *testing.T) {
	pattern := regexp.MustCompile(`^db`)
	opts := classify.DefaultOptions()
	opts.Rules = []classify.Rule{
		{Pattern: pattern, Category: "database", Field: "service"},
	}
	entries := []parser.Entry{
		makeEntry("query slow", map[string]interface{}{"service": "db-primary"}),
		makeEntry("cache miss", map[string]interface{}{"service": "cache"}),
	}
	out := classify.Apply(entries, opts)
	if got := out[0].Fields["category"]; got != "database" {
		t.Errorf("expected 'database', got %v", got)
	}
	if _, ok := out[1].Fields["category"]; ok {
		t.Error("non-matching entry should have no category")
	}
}

func TestApply_DefaultCategory_Applied(t *testing.T) {
	opts := classify.DefaultOptions()
	opts.DefaultCategory = "general"
	entries := []parser.Entry{makeEntry("unrelated message", nil)}
	out := classify.Apply(entries, opts)
	if got := out[0].Fields["category"]; got != "general" {
		t.Errorf("expected 'general', got %v", got)
	}
}

func TestApply_CustomOutputField(t *testing.T) {
	pattern := regexp.MustCompile(`error`)
	opts := classify.DefaultOptions()
	opts.OutputField = "tag"
	opts.Rules = []classify.Rule{
		{Pattern: pattern, Category: "fault"},
	}
	entries := []parser.Entry{makeEntry("disk error detected", nil)}
	out := classify.Apply(entries, opts)
	if got := out[0].Fields["tag"]; got != "fault" {
		t.Errorf("expected 'fault' in 'tag', got %v", got)
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	pattern := regexp.MustCompile(`.*`)
	opts := classify.DefaultOptions()
	opts.Rules = []classify.Rule{
		{Pattern: pattern, Category: "all"},
	}
	original := makeEntry("test", nil)
	classify.Apply([]parser.Entry{original}, opts)
	if _, ok := original.Fields["category"]; ok {
		t.Error("original entry was mutated")
	}
}
