package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParseLine_ValidEntry(t *testing.T) {
	line := "2024-01-15T10:30:00Z INFO server started port=8080"
	entry := parseLine(line)

	if entry.Level != LevelInfo {
		t.Errorf("expected INFO, got %s", entry.Level)
	}
	if entry.Message != "server started" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
	if entry.Fields["port"] != "8080" {
		t.Errorf("expected port=8080, got %q", entry.Fields["port"])
	}
	if entry.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestParseLine_UnknownLevel(t *testing.T) {
	line := "some random log line without structure"
	entry := parseLine(line)

	if entry.Level != LevelUnknown {
		t.Errorf("expected UNKNOWN, got %s", entry.Level)
	}
	if entry.Raw != line {
		t.Errorf("raw should equal original line")
	}
}

func TestParseLine_ErrorLevel(t *testing.T) {
	line := "2024-01-15 09:00:00 ERROR database connection failed retries=3"
	entry := parseLine(line)

	if entry.Level != LevelError {
		t.Errorf("expected ERROR, got %s", entry.Level)
	}
	if entry.Fields["retries"] != "3" {
		t.Errorf("expected retries=3, got %q", entry.Fields["retries"])
	}
}

func TestParser_Parse_MultipleLines(t *testing.T) {
	input := strings.Join([]string{
		"2024-01-15T08:00:00Z DEBUG initializing",
		"2024-01-15T08:00:01Z INFO ready",
		"",
		"2024-01-15T08:00:02Z ERROR crashed",
	}, "\n")

	p := New(strings.NewReader(input))
	entries, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestParseTimestamp_Formats(t *testing.T) {
	cases := []string{
		"2024-01-15T10:30:00Z",
		"2024-01-15T10:30:00",
		"2024-01-15 10:30:00",
	}
	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			ts := parseTimestamp(tc)
			if ts == (time.Time{}) {
				t.Errorf("failed to parse timestamp: %q", tc)
			}
		})
	}
}
