// Package pipeline wires together parsing, filtering, sampling, dedup,
// and output formatting into a single reusable processing chain.
package pipeline

import (
	"io"

	"github.com/yourorg/logslice/internal/dedup"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/sampler"
	"github.com/yourorg/logslice/internal/stats"
)

// Options configures the pipeline.
type Options struct {
	Filter  filter.Options
	Sampler sampler.Options
	Dedup   dedup.Options
	Format  output.Format
	Color   bool
	Stats   bool
}

// Result holds pipeline output metadata.
type Result struct {
	Summary stats.Summary
	Written int
}

// Run reads log lines from r, processes them through the pipeline, and writes
// formatted output to w. It returns a Result with summary statistics.
func Run(r io.Reader, w io.Writer, opts Options) (Result, error) {
	p := parser.New(r)
	entries, err := p.Parse()
	if err != nil {
		return Result{}, fmt.Errorf("pipeline: parsing failed: %w", err)
	}

	filtered := filter.Filter(entries, opts.Filter)
	sampled := sampler.New(opts.Sampler).Apply(filtered)
	deduped := dedup.Dedup(sampled, opts.Dedup)

	fmtr := output.New(opts.Format, opts.Color)
	if err := fmtr.Write(w, deduped); err != nil {
		return Result{}, fmt.Errorf("pipeline: writing output failed: %w", err)
	}

	var summary stats.Summary
	if opts.Stats {
		summary = stats.Compute(deduped)
	}

	return Result{Summary: summary, Written: len(deduped)}, nil
}
