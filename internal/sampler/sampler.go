// Package sampler provides log entry sampling functionality,
// allowing users to reduce output by taking every Nth entry
// or a random percentage of entries.
package sampler

import (
	"math/rand"

	"github.com/user/logslice/internal/parser"
)

// Options controls how sampling is applied.
type Options struct {
	// Every N entries to keep (0 = disabled).
	EveryN int
	// Percent is a value from 0–100; only entries randomly selected
	// within this percentage are kept (0 = disabled).
	Percent int
}

// Sampler holds state for deterministic or probabilistic sampling.
type Sampler struct {
	opts    Options
	counter int
	rng     *rand.Rand
}

// New returns a new Sampler configured with the given Options.
func New(opts Options, seed int64) *Sampler {
	return &Sampler{
		opts: opts,
		rng:  rand.New(rand.NewSource(seed)),
	}
}

// Apply filters entries according to the configured sampling strategy.
// If neither EveryN nor Percent is set, all entries are returned unchanged.
func (s *Sampler) Apply(entries []parser.Entry) []parser.Entry {
	if s.opts.EveryN <= 0 && s.opts.Percent <= 0 {
		return entries
	}

	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		if s.keep(e) {
			out = append(out, e)
		}
	}
	return out
}

// keep returns true if the entry should be included in the sample.
func (s *Sampler) keep(_ parser.Entry) bool {
	if s.opts.EveryN > 0 {
		s.counter++
		if s.counter%s.opts.EveryN != 0 {
			return false
		}
	}
	if s.opts.Percent > 0 {
		if s.rng.Intn(100) >= s.opts.Percent {
			return false
		}
	}
	return true
}
