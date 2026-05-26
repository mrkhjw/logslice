package fieldmap_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/fieldmap"
	"github.com/yourorg/logslice/internal/parser"
)

func makeEntry(fields map[string]string) parser.Entry {
	return parser.Entry{
		Timestamp: time.Now(),
		Level:     "info",
		Message:   "test message",
		Fields:    fields,
	}
}

func TestApply_NoRules_ReturnsOriginal(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{"host": "srv1"})}
	out := fieldmap.Apply(entries, fieldmap.DefaultOptions())
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
	if out[0].Fields["host"] != "srv1" {
		t.Errorf("expected host=srv1, got %q", out[0].Fields["host"])
	}
}

func TestApply_RenamesField(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{"host": "srv1", "svc": "api"})}
	opts := fieldmap.Options{
		Map: []fieldmap.Rule{{Src: "host", Dst: "hostname"}},
	}
	out := fieldmap.Apply(entries, opts)
	if _, ok := out[0].Fields["host"]; ok {
		t.Error("old key 'host' should have been removed")
	}
	if out[0].Fields["hostname"] != "srv1" {
		t.Errorf("expected hostname=srv1, got %q", out[0].Fields["hostname"])
	}
	if out[0].Fields["svc"] != "api" {
		t.Error("unmapped field 'svc' should be preserved")
	}
}

func TestApply_DropUnmapped(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{"host": "srv1", "svc": "api", "env": "prod"})}
	opts := fieldmap.Options{
		Map:          []fieldmap.Rule{{Src: "host", Dst: "hostname"}},
		DropUnmapped: true,
	}
	out := fieldmap.Apply(entries, opts)
	if len(out[0].Fields) != 1 {
		t.Errorf("expected 1 field after drop, got %d: %v", len(out[0].Fields), out[0].Fields)
	}
	if out[0].Fields["hostname"] != "srv1" {
		t.Errorf("expected hostname=srv1, got %q", out[0].Fields["hostname"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	orig := makeEntry(map[string]string{"host": "srv1"})
	opts := fieldmap.Options{
		Map: []fieldmap.Rule{{Src: "host", Dst: "hostname"}},
	}
	fieldmap.Apply([]parser.Entry{orig}, opts)
	if _, ok := orig.Fields["host"]; !ok {
		t.Error("original entry should not be mutated")
	}
}

func TestApply_EmptyInput(t *testing.T) {
	out := fieldmap.Apply(nil, fieldmap.DefaultOptions())
	if len(out) != 0 {
		t.Errorf("expected empty output, got %d entries", len(out))
	}
}

func TestApply_MultipleRules(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]string{"a": "1", "b": "2", "c": "3"})}
	opts := fieldmap.Options{
		Map: []fieldmap.Rule{
			{Src: "a", Dst: "alpha"},
			{Src: "b", Dst: "beta"},
		},
	}
	out := fieldmap.Apply(entries, opts)
	if out[0].Fields["alpha"] != "1" || out[0].Fields["beta"] != "2" || out[0].Fields["c"] != "3" {
		t.Errorf("unexpected fields: %v", out[0].Fields)
	}
}
