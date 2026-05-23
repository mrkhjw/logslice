// Package pipeline provides a high-level processing chain for log entries.
//
// It composes the parser, filter, sampler, dedup, and output packages into a
// single Run function that reads raw log lines, applies all transformations,
// and writes formatted results to an io.Writer.
//
// Usage:
//
//	res, err := pipeline.Run(os.Stdin, os.Stdout, pipeline.Options{
//		Format: output.FormatJSON,
//		Filter: filter.Options{Level: "ERROR"},
//		Stats:  true,
//	})
package pipeline
