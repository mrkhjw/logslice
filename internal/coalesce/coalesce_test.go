package coalesce_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/coalesce"
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
	entries := []parser.Entry{makeEntry(map[string]string{"a": "1"})}
	got := coalesce.Apply(entries, coalesce.DefaultOptions())
	if len(got) != 1 || got[0].Fields["a"] != "1" {
		t.Fatal("expected original entry unchanged")
	}
}

func TestApply_FirstFieldWins(t *testing.T) {
	opts := coalesce.DefaultOptions()
	opts.Rules = []coalesce.Rule{{Target: "host", Fields: []string{"hostname", "node", "server"}}}

	e := makeEntry(map[string]string{"node": "node-1", "server": "srv-2"})
	got := coalesce.Apply([]parser.Entry{e}, opts)

	if got[0].Fields["host"] != "node-1" {
		t.Fatalf("expected node-1, got %q", got[0].Fields["host"])
	}
}

func TestApply_SkipsEmptyFields(t *testing.T) {
	opts := coalesce.DefaultOptions()
	opts.Rules = []coalesce.Rule{{Target: "svc", Fields: []string{"service", "app"}}}

	e := makeEntry(map[string]string{"service": "", "app": "myapp"})
	got := coalesce.Apply([]parser.Entry{e}, opts)

	if got[0].Fields["svc"] != "myapp" {
		t.Fatalf("expected myapp, got %q", got[0].Fields["svc"])
	}
}

func TestApply_NoMatchLeavesMissing(t *testing.T) {
	opts := coalesce.DefaultOptions()
	opts.Rules = []coalesce.Rule{{Target: "host", Fields: []string{"hostname"}}}

	e := makeEntry(map[string]string{"other": "val"})
	got := coalesce.Apply([]parser.Entry{e}, opts)

	if _, ok := got[0].Fields["host"]; ok {
		t.Fatal("expected host field to be absent")
	}
}

func TestApply_DropSources(t *testing.T) {
	opts := coalesce.Options{
		DropSources: true,
		Rules:       []coalesce.Rule{{Target: "host", Fields: []string{"hostname", "node"}}},
	}

	e := makeEntry(map[string]string{"hostname": "h1", "node": "n1"})
	got := coalesce.Apply([]parser.Entry{e}, opts)

	if got[0].Fields["host"] != "h1" {
		t.Fatalf("expected h1, got %q", got[0].Fields["host"])
	}
	if _, ok := got[0].Fields["hostname"]; ok {
		t.Fatal("expected hostname to be dropped")
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	opts := coalesce.DefaultOptions()
	opts.Rules = []coalesce.Rule{{Target: "host", Fields: []string{"hostname"}}}

	origFields := map[string]string{"hostname": "h1"}
	e := makeEntry(origFields)
	coalesce.Apply([]parser.Entry{e}, opts)

	if _, ok := e.Fields["host"]; ok {
		t.Fatal("original entry must not be mutated")
	}
}
