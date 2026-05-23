package truncate_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/truncate"
)

func makeEntry(msg string, fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   msg,
		Fields:    fields,
	}
}

func TestApply_ShortMessage_Unchanged(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello", nil)}
	opts := truncate.DefaultOptions()
	out := truncate.Apply(entries, opts)
	if out[0].Message != "hello" {
		t.Errorf("expected 'hello', got %q", out[0].Message)
	}
}

func TestApply_LongMessage_Truncated(t *testing.T) {
	long := strings.Repeat("a", 300)
	entries := []parser.Entry{makeEntry(long, nil)}
	opts := truncate.DefaultOptions()
	out := truncate.Apply(entries, opts)
	if len([]rune(out[0].Message)) != truncate.DefaultMaxLength+len([]rune(truncate.DefaultSuffix)) {
		t.Errorf("unexpected length: %d", len(out[0].Message))
	}
	if !strings.HasSuffix(out[0].Message, truncate.DefaultSuffix) {
		t.Errorf("expected suffix %q", truncate.DefaultSuffix)
	}
}

func TestApply_CustomMaxLength(t *testing.T) {
	entries := []parser.Entry{makeEntry("hello world", nil)}
	opts := truncate.Options{MaxLength: 5, Suffix: "~", TruncateFields: false}
	out := truncate.Apply(entries, opts)
	if out[0].Message != "hello~" {
		t.Errorf("expected 'hello~', got %q", out[0].Message)
	}
}

func TestApply_TruncatesFields(t *testing.T) {
	fields := map[string]string{"key": strings.Repeat("x", 300)}
	entries := []parser.Entry{makeEntry("msg", fields)}
	opts := truncate.DefaultOptions()
	out := truncate.Apply(entries, opts)
	v := out[0].Fields["key"]
	if !strings.HasSuffix(v, truncate.DefaultSuffix) {
		t.Errorf("field value not truncated: %q", v)
	}
}

func TestApply_DoesNotTruncateFields_WhenDisabled(t *testing.T) {
	long := strings.Repeat("y", 300)
	fields := map[string]string{"k": long}
	entries := []parser.Entry{makeEntry("msg", fields)}
	opts := truncate.Options{MaxLength: 10, Suffix: "...", TruncateFields: false}
	out := truncate.Apply(entries, opts)
	if out[0].Fields["k"] != long {
		t.Errorf("field should not be truncated")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	origMsg := strings.Repeat("z", 300)
	entries := []parser.Entry{makeEntry(origMsg, nil)}
	opts := truncate.DefaultOptions()
	truncate.Apply(entries, opts)
	if entries[0].Message != origMsg {
		t.Errorf("original entry was mutated")
	}
}

func TestDefaultOptions_Values(t *testing.T) {
	opts := truncate.DefaultOptions()
	if opts.MaxLength != truncate.DefaultMaxLength {
		t.Errorf("expected MaxLength %d, got %d", truncate.DefaultMaxLength, opts.MaxLength)
	}
	if opts.Suffix != truncate.DefaultSuffix {
		t.Errorf("expected Suffix %q, got %q", truncate.DefaultSuffix, opts.Suffix)
	}
	if !opts.TruncateFields {
		t.Error("expected TruncateFields to be true")
	}
}
