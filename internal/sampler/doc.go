// Package sampler provides log entry sampling strategies for logslice.
//
// Two strategies are supported:
//
//   - EveryN: keep every Nth log entry in sequence (deterministic).
//   - Percent: keep a random percentage of entries (probabilistic).
//
// Both strategies can be combined; an entry must satisfy both conditions
// to be included in the output.
//
// Example usage:
//
//	s := sampler.New(sampler.Options{EveryN: 10, Percent: 80}, time.Now().UnixNano())
//	filtered := s.Apply(entries)
package sampler
