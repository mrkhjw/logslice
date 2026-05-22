package output_test

import (
	"testing"

	"github.com/user/logslice/internal/output"
)

func TestParseFormat_Valid(t *testing.T) {
	cases := []struct {
		input    string
		expected output.Format
	}{
		{"json", output.FormatJSON},
		{"JSON", output.FormatJSON},
		{"text", output.FormatText},
		{"TEXT", output.FormatText},
		{"", output.FormatText},
		{"table", output.FormatTable},
		{"TABLE", output.FormatTable},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got, err := output.ParseFormat(tc.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.expected {
				t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestParseFormat_Invalid(t *testing.T) {
	_, err := output.ParseFormat("xml")
	if err == nil {
		t.Error("expected error for unknown format, got nil")
	}
}

func TestValidFormats(t *testing.T) {
	formats := output.ValidFormats()
	if len(formats) != 3 {
		t.Errorf("expected 3 formats, got %d", len(formats))
	}
}

func TestFormat_String(t *testing.T) {
	if output.FormatJSON.String() != "json" {
		t.Errorf("expected \"json\", got %q", output.FormatJSON.String())
	}
}
