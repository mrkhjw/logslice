package pipeline_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/pipeline"
	"github.com/yourorg/logslice/internal/sampler"
)

func TestBuilder_Defaults(t *testing.T) {
	opts, err := pipeline.NewBuilder().Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != output.FormatText {
		t.Errorf("expected default format text, got %v", opts.Format)
	}
	if !opts.Color {
		t.Error("expected color enabled by default")
	}
}

func TestBuilder_WithLevel(t *testing.T) {
	opts, err := pipeline.NewBuilder().WithLevel("WARN").Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Filter.Level != "WARN" {
		t.Errorf("expected level WARN, got %s", opts.Filter.Level)
	}
}

func TestBuilder_WithTimeRange(t *testing.T) {
	since := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	until := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

	opts, err := pipeline.NewBuilder().WithTimeRange(since, until).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Filter.Since.Equal(since) {
		t.Errorf("since mismatch")
	}
	if !opts.Filter.Until.Equal(until) {
		t.Errorf("until mismatch")
	}
}

func TestBuilder_WithSampler(t *testing.T) {
	s := sampler.Options{EveryN: 2}
	opts, err := pipeline.NewBuilder().WithSampler(s).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Sampler.EveryN != 2 {
		t.Errorf("expected EveryN 2, got %d", opts.Sampler.EveryN)
	}
}

func TestBuilder_WithStats(t *testing.T) {
	opts, err := pipeline.NewBuilder().WithStats(true).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Stats {
		t.Error("expected stats enabled")
	}
}

func TestBuilder_WithFields(t *testing.T) {
	fields := map[string]string{"service": "api", "env": "prod"}
	opts, err := pipeline.NewBuilder().WithFields(fields).Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Filter.Fields["service"] != "api" {
		t.Errorf("expected service=api in fields")
	}
}
