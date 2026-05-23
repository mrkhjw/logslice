package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
)

const sampleLog = `2024-01-15T10:00:00Z INFO  service started
2024-01-15T10:00:01Z DEBUG loading config
2024-01-15T10:00:02Z ERROR connection refused
2024-01-15T10:00:03Z INFO  ready
2024-01-15T10:00:04Z WARN  high memory usage
`

func TestRun_AllEntries(t *testing.T) {
	r := strings.NewReader(sampleLog)
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Options{Format: output.FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 5 {
		t.Errorf("expected 5 written, got %d", res.Written)
	}
}

func TestRun_FilterByLevel(t *testing.T) {
	r := strings.NewReader(sampleLog)
	var w bytes.Buffer

	opts := pipeline.Options{
		Format: output.FormatText,
		Filter: filter.Options{Level: "ERROR"},
	}
	res, err := pipeline.Run(r, &w, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 1 {
		t.Errorf("expected 1 written, got %d", res.Written)
	}
	if !strings.Contains(w.String(), "connection refused") {
		t.Errorf("expected error entry in output")
	}
}

func TestRun_JSONFormat(t *testing.T) {
	r := strings.NewReader(sampleLog)
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Options{Format: output.FormatJSON})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written == 0 {
		t.Error("expected entries to be written")
	}
	if !strings.Contains(w.String(), `"level"`) {
		t.Errorf("expected JSON output with level field")
	}
}

func TestRun_WithStats(t *testing.T) {
	r := strings.NewReader(sampleLog)
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Options{
		Format: output.FormatText,
		Stats:  true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Summary.Total != 5 {
		t.Errorf("expected summary total 5, got %d", res.Summary.Total)
	}
}

func TestRun_EmptyInput(t *testing.T) {
	r := strings.NewReader("")
	var w bytes.Buffer

	res, err := pipeline.Run(r, &w, pipeline.Options{Format: output.FormatText})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Written != 0 {
		t.Errorf("expected 0 written, got %d", res.Written)
	}
}
