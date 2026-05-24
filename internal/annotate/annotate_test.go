package annotate_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/annotate"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	if fields == nil {
		fields = map[string]string{}
	}
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello world", nil)}
	result := annotate.Apply(entries, annotate.DefaultOptions())
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0].Message != "hello world" {
		t.Errorf("unexpected message: %s", result[0].Message)
	}
}

func TestApply_MatchMessage_AddsField(t *testing.T) {
	entry := makeEntry("user=alice logged in", nil)
	opts := annotate.Options{
		Rules: []annotate.Rule{
			{
				Field:   "message",
				Pattern: regexp.MustCompile(`user=(\w+)`),
				Key:     "username",
				Value:   "$1",
			},
		},
	}
	result := annotate.Apply([]parser.Entry{entry}, opts)
	if got := result[0].Fields["username"]; got != "alice" {
		t.Errorf("expected username=alice, got %q", got)
	}
}

func TestApply_NoMatch_FieldAbsent(t *testing.T) {
	entry := makeEntry("nothing interesting", nil)
	opts := annotate.Options{
		Rules: []annotate.Rule{
			{
				Pattern: regexp.MustCompile(`user=(\w+)`),
				Key:     "username",
				Value:   "$1",
			},
		},
	}
	result := annotate.Apply([]parser.Entry{entry}, opts)
	if _, ok := result[0].Fields["username"]; ok {
		t.Error("expected username field to be absent")
	}
}

func TestApply_MatchNamedField(t *testing.T) {
	entry := makeEntry("event fired", map[string]string{"env": "prod-us-east"})
	opts := annotate.Options{
		Rules: []annotate.Rule{
			{
				Field:   "env",
				Pattern: regexp.MustCompile(`^(\w+)-`),
				Key:     "env_short",
				Value:   "$1",
			},
		},
	}
	result := annotate.Apply([]parser.Entry{entry}, opts)
	if got := result[0].Fields["env_short"]; got != "prod" {
		t.Errorf("expected env_short=prod, got %q", got)
	}
}

func TestApply_StopOnFirst(t *testing.T) {
	entry := makeEntry("error: disk full", nil)
	opts := annotate.Options{
		StopOnFirst: true,
		Rules: []annotate.Rule{
			{
				Pattern: regexp.MustCompile(`error`),
				Key:     "tag",
				Value:   "matched",
			},
			{
				Pattern: regexp.MustCompile(`disk`),
				Key:     "second_tag",
				Value:   "also_matched",
			},
		},
	}
	result := annotate.Apply([]parser.Entry{entry}, opts)
	if result[0].Fields["tag"] != "matched" {
		t.Error("expected first rule to match")
	}
	if _, ok := result[0].Fields["second_tag"]; ok {
		t.Error("expected second rule to be skipped due to StopOnFirst")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	original := makeEntry("hello", nil)
	opts := annotate.Options{
		Rules: []annotate.Rule{
			{
				Pattern: regexp.MustCompile(`hello`),
				Key:     "greeted",
				Value:   "true",
			},
		},
	}
	annotate.Apply([]parser.Entry{original}, opts)
	if _, ok := original.Fields["greeted"]; ok {
		t.Error("original entry was mutated")
	}
}
