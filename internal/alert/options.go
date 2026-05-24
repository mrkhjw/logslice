package alert

import "time"

// Builder provides a fluent API for constructing a set of alert Rules.
type Builder struct {
	rules []Rule
}

// NewBuilder returns an empty Builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// OnLevel adds a rule that fires when the count of entries at level
// exceeds threshold within the given window duration.
func (b *Builder) OnLevel(level string, threshold int, window time.Duration) *Builder {
	b.rules = append(b.rules, Rule{
		Level:     level,
		Threshold: threshold,
		Window:    window,
	})
	return b
}

// OnError is a convenience method for adding an error-level rule.
func (b *Builder) OnError(threshold int, window time.Duration) *Builder {
	return b.OnLevel("error", threshold, window)
}

// OnFatal is a convenience method for adding a fatal-level rule.
func (b *Builder) OnFatal(threshold int, window time.Duration) *Builder {
	return b.OnLevel("fatal", threshold, window)
}

// Build returns the constructed slice of Rules.
func (b *Builder) Build() []Rule {
	out := make([]Rule, len(b.rules))
	copy(out, b.rules)
	return out
}
