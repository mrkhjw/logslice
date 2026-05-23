package grep_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/grep"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_MatchMessage(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("user login succeeded", nil),
		makeEntry("disk full error", nil),
		makeEntry("user logout", nil),
	}
	c, err := grep.Compile(grep.Options{Pattern: "user"})
	if err != nil {
		t.Fatal(err)
	}
	got := c.Apply(entries)
	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}

func TestApply_NoMatch_ReturnsEmpty(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello world", nil)}
	c, _ := grep.Compile(grep.Options{Pattern: "notfound"})
	got := c.Apply(entries)
	if len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestApply_Invert(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("error occurred", nil),
		makeEntry("all systems nominal", nil),
	}
	c, _ := grep.Compile(grep.Options{Pattern: "error", Invert: true})
	got := c.Apply(entries)
	if len(got) != 1 || got[0].Message != "all systems nominal" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestApply_FieldPattern(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("msg", map[string]string{"service": "auth"}),
		makeEntry("msg", map[string]string{"service": "billing"}),
	}
	c, _ := grep.Compile(grep.Options{FieldPattern: map[string]string{"service": "^auth$"}})
	got := c.Apply(entries)
	if len(got) != 1 || got[0].Fields["service"] != "auth" {
		t.Fatalf("unexpected result: %+v", got)
	}
}

func TestApply_FieldMissing_Excluded(t *testing.T) {
	entries := []parser.Entry{
		makeEntry("no field", nil),
		makeEntry("has field", map[string]string{"env": "prod"}),
	}
	c, _ := grep.Compile(grep.Options{FieldPattern: map[string]string{"env": "prod"}})
	got := c.Apply(entries)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
}

func TestCompile_InvalidPattern(t *testing.T) {
	_, err := grep.Compile(grep.Options{Pattern: "[invalid"})
	if err == nil {
		t.Fatal("expected compile error for invalid regex")
	}
}

func TestCompile_InvalidFieldPattern(t *testing.T) {
	_, err := grep.Compile(grep.Options{FieldPattern: map[string]string{"field": "[bad"}})
	if err == nil {
		t.Fatal("expected compile error for invalid field regex")
	}
}
