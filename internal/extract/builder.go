package extract

// Builder provides a fluent API for constructing extract Options.
type Builder struct {
	opts Options
}

// NewBuilder returns a Builder initialised with DefaultOptions.
func NewBuilder() *Builder {
	return &Builder{opts: DefaultOptions()}
}

// Field schedules extraction of src, keeping the same name in the output.
func (b *Builder) Field(src string) *Builder {
	b.opts.Fields[src] = ""
	return b
}

// FieldAs schedules extraction of src, renaming it to dst in the output.
func (b *Builder) FieldAs(src, dst string) *Builder {
	b.opts.Fields[src] = dst
	return b
}

// WithMessage controls whether the message is included (default: true).
func (b *Builder) WithMessage(v bool) *Builder {
	b.opts.IncludeMessage = v
	return b
}

// WithLevel controls whether the level is included (default: true).
func (b *Builder) WithLevel(v bool) *Builder {
	b.opts.IncludeLevel = v
	return b
}

// WithTimestamp controls whether the timestamp is included (default: true).
func (b *Builder) WithTimestamp(v bool) *Builder {
	b.opts.IncludeTimestamp = v
	return b
}

// Build returns the constructed Options.
func (b *Builder) Build() Options {
	return b.opts
}
