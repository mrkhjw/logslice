// Package classify assigns a category tag to log entries based on
// user-defined pattern rules.
//
// Each rule specifies a regular expression and a category string. Rules are
// evaluated in order; the first match wins. Patterns can be applied against
// the entry message or any named field.
//
// Basic usage:
//
//	import (
//		"regexp"
//		"github.com/yourorg/logslice/internal/classify"
//	)
//
//	opts := classify.DefaultOptions()
//	opts.Rules = []classify.Rule{
//		{Pattern: regexp.MustCompile(`timeout`), Category: "network"},
//		{Pattern: regexp.MustCompile(`panic`),   Category: "crash"},
//	}
//	opts.DefaultCategory = "general"
//
//	annotated := classify.Apply(entries, opts)
//
// The resulting entries carry a "category" field (configurable via
// Options.OutputField) that downstream components such as route or aggregate
// can use for grouping.
package classify
