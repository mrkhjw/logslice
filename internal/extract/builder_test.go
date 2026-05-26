package extract_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/extract"
)

func TestBuilder_Defaults(t *testing.T) {
	opts := extract.NewBuilder().Build()
	if !opts.IncludeMessage {
		t.Error("expected IncludeMessage=true by default")
	}
	if !opts.IncludeLevel {
		t.Error("expected IncludeLevel=true by default")
	}
	if !opts.IncludeTimestamp {
		t.Error("expected IncludeTimestamp=true by default")
	}
	if len(opts.Fields) != 0 {
		t.Errorf("expected no fields by default, got %v", opts.Fields)
	}
}

func TestBuilder_Field_KeepsName(t *testing.T) {
	opts := extract.NewBuilder().Field("user").Build()
	dst, ok := opts.Fields["user"]
	if !ok {
		t.Fatal("expected 'user' field to be registered")
	}
	if dst != "" {
		t.Errorf("expected empty dst (keep name), got %q", dst)
	}
}

func TestBuilder_FieldAs_RenamesField(t *testing.T) {
	opts := extract.NewBuilder().FieldAs("usr", "user").Build()
	dst, ok := opts.Fields["usr"]
	if !ok {
		t.Fatal("expected 'usr' field to be registered")
	}
	if dst != "user" {
		t.Errorf("expected dst=user, got %q", dst)
	}
}

func TestBuilder_WithMessage_False(t *testing.T) {
	opts := extract.NewBuilder().WithMessage(false).Build()
	if opts.IncludeMessage {
		t.Error("expected IncludeMessage=false")
	}
}

func TestBuilder_WithLevel_False(t *testing.T) {
	opts := extract.NewBuilder().WithLevel(false).Build()
	if opts.IncludeLevel {
		t.Error("expected IncludeLevel=false")
	}
}

func TestBuilder_WithTimestamp_False(t *testing.T) {
	opts := extract.NewBuilder().WithTimestamp(false).Build()
	if opts.IncludeTimestamp {
		t.Error("expected IncludeTimestamp=false")
	}
}

func TestBuilder_Chaining(t *testing.T) {
	opts := extract.NewBuilder().
		Field("host").
		FieldAs("svc", "service").
		WithMessage(false).
		Build()

	if opts.IncludeMessage {
		t.Error("expected IncludeMessage=false")
	}
	if len(opts.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(opts.Fields))
	}
}
