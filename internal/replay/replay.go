// Package replay provides functionality to replay log entries at their
// original speed or at a scaled rate, useful for testing and simulation.
package replay

import (
	"time"

	"github.com/yourorg/logslice/internal/parser"
)

// Options controls replay behaviour.
type Options struct {
	// Speed is a multiplier for the original inter-entry delay.
	// 1.0 = real time, 2.0 = twice as fast, 0.5 = half speed.
	// Values <= 0 default to 1.0.
	Speed float64

	// MaxDelay caps the pause between entries so a long silence in the
	// original log does not stall the replay indefinitely.
	MaxDelay time.Duration
}

// DefaultOptions returns sensible defaults for replay.
func DefaultOptions() Options {
	return Options{
		Speed:    1.0,
		MaxDelay: 5 * time.Second,
	}
}

// Apply sends entries to the returned channel, pausing between each pair
// of entries to simulate the original timing. The channel is closed when
// all entries have been sent.
func Apply(entries []parser.Entry, opts Options) <-chan parser.Entry {
	if opts.Speed <= 0 {
		opts.Speed = 1.0
	}
	if opts.MaxDelay <= 0 {
		opts.MaxDelay = DefaultOptions().MaxDelay
	}

	ch := make(chan parser.Entry, 1)

	go func() {
		defer close(ch)
		for i, entry := range entries {
			if i > 0 {
				prev := entries[i-1]
				if !entry.Timestamp.IsZero() && !prev.Timestamp.IsZero() {
					gap := entry.Timestamp.Sub(prev.Timestamp)
					if gap > 0 {
						scaled := time.Duration(float64(gap) / opts.Speed)
						if scaled > opts.MaxDelay {
							scaled = opts.MaxDelay
						}
						time.Sleep(scaled)
					}
				}
			}
			ch <- entry
		}
	}()

	return ch
}
