package pipeline

import (
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/dedup"
	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/output"
	"github.com/yourorg/logslice/internal/sampler"
)

// Builder helps construct Options with a fluent API.
type Builder struct {
	opts Options
}

// NewBuilder returns a Builder with sensible defaults.
func NewBuilder() *Builder {
	return &Builder{
		opts: Options{
			Format: output.FormatText,
			Color:  true,
		},
	}
}

// WithFormat sets the output format.
func (b *Builder) WithFormat(f output.Format) *Builder {
	b.opts.Format = f
	return b
}

// WithLevel restricts output to entries matching level.
func (b *Builder) WithLevel(level string) *Builder {
	b.opts.Filter.Level = level
	return b
}

// WithTimeRange restricts output to entries within [since, until].
func (b *Builder) WithTimeRange(since, until time.Time) *Builder {
	b.opts.Filter.Since = since
	b.opts.Filter.Until = until
	return b
}

// WithSampler configures entry sampling.
func (b *Builder) WithSampler(s sampler.Options) *Builder {
	b.opts.Sampler = s
	return b
}

// WithDedup configures deduplication.
func (b *Builder) WithDedup(d dedup.Options) *Builder {
	b.opts.Dedup = d
	return b
}

// WithStats enables statistics collection.
func (b *Builder) WithStats(enabled bool) *Builder {
	b.opts.Stats = enabled
	return b
}

// WithColor enables or disables ANSI colour output.
func (b *Builder) WithColor(enabled bool) *Builder {
	b.opts.Color = enabled
	return b
}

// WithFields adds key=value field filters.
func (b *Builder) WithFields(fields map[string]string) *Builder {
	b.opts.Filter.Fields = fields
	return b
}

// Build validates and returns the configured Options.
func (b *Builder) Build() (Options, error) {
	if b.opts.Format < 0 {
		return Options{}, fmt.Errorf("pipeline: invalid output format %d", b.opts.Format)
	}
	return b.opts, nil
}
