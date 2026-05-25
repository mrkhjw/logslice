package route_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/route"
)

func makeEntry(level parser.Level, msg string, fields map[string]string) parser.Entry {
	e := parser.Entry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   msg,
		Fields:    make(map[string]string),
	}
	for k, v := range fields {
		e.Fields[k] = v
	}
	return e
}

func TestApply_NoRules_AllUnmatched(t *testing.T) {
	entries := []parser.Entry{makeEntry(parser.LevelInfo, "hello", nil)}
	res := route.Apply(entries, nil)
	if len(res.Unmatched) != 1 {
		t.Fatalf("expected 1 unmatched, got %d", len(res.Unmatched))
	}
	if len(res.Buckets) != 0 {
		t.Fatalf("expected no buckets, got %d", len(res.Buckets))
	}
}

func TestApply_RouteByLevel(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelError, "boom", nil),
		makeEntry(parser.LevelInfo, "ok", nil),
	}
	rules := []route.Rule{
		{Field: "level", Pattern: regexp.MustCompile(`(?i)error`), Dest: "errors"},
	}
	res := route.Apply(entries, rules)
	if len(res.Buckets["errors"]) != 1 {
		t.Fatalf("expected 1 error entry, got %d", len(res.Buckets["errors"]))
	}
	if len(res.Unmatched) != 1 {
		t.Fatalf("expected 1 unmatched, got %d", len(res.Unmatched))
	}
}

func TestApply_RouteByMessage(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, "user login", nil),
		makeEntry(parser.LevelInfo, "payment processed", nil),
	}
	rules := []route.Rule{
		{Field: "message", Pattern: regexp.MustCompile(`payment`), Dest: "billing"},
	}
	res := route.Apply(entries, rules)
	if len(res.Buckets["billing"]) != 1 {
		t.Fatalf("expected 1 billing entry, got %d", len(res.Buckets["billing"]))
	}
}

func TestApply_RouteByCustomField(t *testing.T) {
	entries := []parser.Entry{
		makeEntry(parser.LevelInfo, "req", map[string]string{"service": "auth"}),
		makeEntry(parser.LevelInfo, "req", map[string]string{"service": "billing"}),
	}
	rules := []route.Rule{
		{Field: "service", Pattern: regexp.MustCompile(`auth`), Dest: "auth-logs"},
	}
	res := route.Apply(entries, rules)
	if len(res.Buckets["auth-logs"]) != 1 {
		t.Fatalf("expected 1 auth-logs entry, got %d", len(res.Buckets["auth-logs"]))
	}
	if len(res.Unmatched) != 1 {
		t.Fatalf("expected 1 unmatched, got %d", len(res.Unmatched))
	}
}

func TestApply_FirstRuleWins(t *testing.T) {
	entries := []parser.Entry{makeEntry(parser.LevelError, "critical failure", nil)}
	rules := []route.Rule{
		{Field: "level", Pattern: regexp.MustCompile(`(?i)error`), Dest: "errors"},
		{Field: "message", Pattern: regexp.MustCompile(`critical`), Dest: "critical"},
	}
	res := route.Apply(entries, rules)
	if len(res.Buckets["errors"]) != 1 {
		t.Fatalf("expected entry in 'errors' bucket")
	}
	if len(res.Buckets["critical"]) != 0 {
		t.Fatalf("expected 'critical' bucket to be empty")
	}
}
