package transform_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/transform"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     parser.LevelInfo,
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_NoOps_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello", nil)}
	result := transform.Apply(entries, nil)
	if len(result) != 1 || result[0].Message != "hello" {
		t.Fatalf("expected unchanged entry, got %+v", result)
	}
}

func TestApply_UpperCaseMessage(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello world", nil)}
	ops := []transform.Op{{Field: "message", UpperCase: true}}
	result := transform.Apply(entries, ops)
	if result[0].Message != "HELLO WORLD" {
		t.Errorf("expected HELLO WORLD, got %q", result[0].Message)
	}
}

func TestApply_LowerCaseField(t *testing.T) {
	entries := []parser.Entry{makeEntry("msg", map[string]string{"service": "MyService"})}
	ops := []transform.Op{{Field: "service", LowerCase: true}}
	result := transform.Apply(entries, ops)
	if result[0].Fields["service"] != "myservice" {
		t.Errorf("expected myservice, got %q", result[0].Fields["service"])
	}
}

func TestApply_RenameField(t *testing.T) {
	entries := []parser.Entry{makeEntry("msg", map[string]string{"svc": "auth"})}
	ops := []transform.Op{{Field: "svc", Rename: "service"}}
	result := transform.Apply(entries, ops)
	if _, ok := result[0].Fields["svc"]; ok {
		t.Error("old field key should have been removed")
	}
	if result[0].Fields["service"] != "auth" {
		t.Errorf("expected renamed field value 'auth', got %q", result[0].Fields["service"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	orig := makeEntry("original", map[string]string{"k": "v"})
	entries := []parser.Entry{orig}
	ops := []transform.Op{{Field: "message", UpperCase: true}, {Field: "k", Rename: "key"}}
	transform.Apply(entries, ops)
	if entries[0].Message != "original" {
		t.Error("original entry message was mutated")
	}
	if _, ok := entries[0].Fields["k"]; !ok {
		t.Error("original entry fields were mutated")
	}
}

func TestApply_MissingField_Skipped(t *testing.T) {
	entries := []parser.Entry{makeEntry("msg", map[string]string{})}
	ops := []transform.Op{{Field: "nonexistent", UpperCase: true}}
	result := transform.Apply(entries, ops)
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
}
