// Package batch provides utilities for grouping log entries into fixed-size
// or time-bounded batches.
//
// # Overview
//
// Batching is useful when downstream consumers (e.g. bulk indexers, alerting
// systems, or file writers) operate more efficiently on groups of entries
// rather than individual ones.
//
// # Usage
//
//	opts := batch.Options{
//	    Size:     500,
//	    Duration: 5 * time.Second,
//	}
//	batches := batch.Apply(entries, opts)
//	for _, b := range batches {
//	    process(b)
//	}
//
// When both Size and Duration are set the batch is flushed as soon as either
// threshold is reached.
package batch
