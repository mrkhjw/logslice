package label

import "regexp"

// Builder provides a fluent API for constructing Options.
type Builder struct {
	opts Options
}

// NewBuilder returns a Builder initialised with DefaultOptions.
func NewBuilder() *Builder {
	return &Builder{opts: DefaultOptions()}
}

// OnMessage adds a rule that matches against the entry message.
func (b *Builder) OnMessage(pattern, lbl string) *Builder {
	b.opts.Rules = append(b.opts.Rules, Rule{
		Pattern: regexp.MustCompile(pattern),
		Label:   lbl,
	})
	return b
}

// OnField adds a rule that matches against a named field value.
func (b *Builder) OnField(field, pattern, lbl string) *Builder {
	b.opts.Rules = append(b.opts.Rules, Rule{
		Field:   field,
		Pattern: regexp.MustCompile(pattern),
		Label:   lbl,
	})
	return b
}

// WithLabelField overrides the destination field name (default: "label").
func (b *Builder) WithLabelField(name string) *Builder {
	b.opts.LabelField = name
	return b
}

// WithMulti enables multi-label mode: all matching rules are applied.
func (b *Builder) WithMulti() *Builder {
	b.opts.Multi = true
	return b
}

// Build returns the constructed Options.
func (b *Builder) Build() Options {
	return b.opts
}
