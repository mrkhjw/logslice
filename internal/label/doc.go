// Package label attaches string labels to log entries based on pattern-matching
// rules. Labels are written into a configurable entry field (default: "label")
// and can later be used for filtering, routing, or display.
//
// # Basic usage
//
//	opts := label.NewBuilder().
//		OnMessage(`timeout`, "slow").
//		OnField("service", `^auth`, "auth-svc").
//		WithMulti().
//		Build()
//
//	labeled := label.Apply(entries, opts)
//
// # Rules
//
// Each Rule holds a compiled regexp, an optional field name, and a label
// string. When Field is empty the pattern is matched against the entry
// Message; otherwise it is matched against the value of the named field.
//
// # Multi-label mode
//
// By default only the first matching rule is applied. Set Options.Multi = true
// (or call Builder.WithMulti) to allow multiple labels to be joined with a
// comma in the destination field.
package label
