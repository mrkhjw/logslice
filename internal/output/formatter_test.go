package output_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/output"
	"github.com/user/logslice/internal/parser"
)

func makeEntry(ts time.Time, level, msg string) parser.Entry {
	return parser.Entry{
		Timestamp: ts,
		Level:     level,
		Message:   msg,
		Fields:    map[string]string{},
	}
}

var testTime = time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

func TestFormatter_Text(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatText)
	entries := []parser.Entry{makeEntry(testTime, "INFO", "server started")}

	if err := f.Write(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "INFO") || !strings.Contains(got, "server started") {
		t.Errorf("unexpected text output: %q", got)
	}
}

func TestFormatter_JSON(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	entries := []parser.Entry{makeEntry(testTime, "ERROR", "disk full")}

	if err := f.Write(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "ERROR") || !strings.Contains(got, "disk full") {
		t.Errorf("unexpected json output: %q", got)
	}
}

func TestFormatter_Table(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatTable)
	entries := []parser.Entry{
		makeEntry(testTime, "WARN", "high memory"),
		makeEntry(testTime.Add(time.Minute), "INFO", "gc ran"),
	}

	if err := f.Write(entries); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "TIMESTAMP") {
		t.Errorf("table missing header: %q", got)
	}
	if !strings.Contains(got, "WARN") || !strings.Contains(got, "high memory") {
		t.Errorf("table missing entry data: %q", got)
	}
}

func TestFormatter_EmptyEntries(t *testing.T) {
	var buf bytes.Buffer
	f := output.New(&buf, output.FormatJSON)
	if err := f.Write([]parser.Entry{}); err != nil {
		t.Fatalf("unexpected error on empty entries: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output, got %q", buf.String())
	}
}
