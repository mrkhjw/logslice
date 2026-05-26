// Package fieldmap remaps, renames, and filters fields in log entries.
//
// It is useful when normalising log entries from disparate sources that use
// different field naming conventions before passing them to downstream
// pipeline stages.
//
// Basic usage:
//
//	opts := fieldmap.Options{
//		Map: []fieldmap.Rule{
//			{Src: "host",    Dst: "hostname"},
//			{Src: "service", Dst: "svc"},
//		},
//		DropUnmapped: false,
//	}
//	out := fieldmap.Apply(entries, opts)
//
// Set DropUnmapped to true to act as an allowlist, keeping only fields
// explicitly listed in the Map rules.
package fieldmap
