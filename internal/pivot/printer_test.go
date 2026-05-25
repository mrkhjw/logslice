package pivot_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/pivot"
)

func TestPrintSummary_ContainsHeader(t *testing.T) {
	res := pivot.Result{
		KeyField: "level",
		Rows: []pivot.Row{
			{Key: "info", Count: 3},
			{Key: "error", Count: 1},
		},
	}
	var sb strings.Builder
	pivot.PrintSummary(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "level") {
		t.Errorf("expected header to contain 'level', got:\n%s", out)
	}
	if !strings.Contains(out, "count") {
		t.Errorf("expected header to contain 'count', got:\n%s", out)
	}
}

func TestPrintSummary_ContainsRows(t *testing.T) {
	res := pivot.Result{
		KeyField: "level",
		Rows: []pivot.Row{
			{Key: "warn", Count: 5},
		},
	}
	var sb strings.Builder
	pivot.PrintSummary(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "warn") {
		t.Errorf("expected 'warn' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "5") {
		t.Errorf("expected count '5' in output, got:\n%s", out)
	}
}

func TestPrintSummary_ValueField(t *testing.T) {
	res := pivot.Result{
		KeyField:   "level",
		ValueField: "host",
		Rows: []pivot.Row{
			{Key: "info", Count: 2, Values: []string{"web1", "web2"}},
		},
	}
	var sb strings.Builder
	pivot.PrintSummary(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "host") {
		t.Errorf("expected 'host' in output, got:\n%s", out)
	}
	if !strings.Contains(out, "web1") {
		t.Errorf("expected 'web1' in output, got:\n%s", out)
	}
}

func TestPrintSummary_EmptyResult(t *testing.T) {
	res := pivot.Result{KeyField: "level"}
	var sb strings.Builder
	pivot.PrintSummary(&sb, res)
	out := sb.String()
	if !strings.Contains(out, "no entries") {
		t.Errorf("expected '(no entries)' message, got:\n%s", out)
	}
}
