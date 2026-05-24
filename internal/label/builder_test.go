package label_test

import (
	"testing"

	"github.com/yourorg/logslice/internal/label"
)

func TestBuilder_Defaults(t *testing.T) {
	opts := label.NewBuilder().Build()
	if opts.LabelField != "label" {
		t.Fatalf("expected default LabelField 'label', got %q", opts.LabelField)
	}
	if opts.Multi {
		t.Fatal("expected Multi to be false by default")
	}
	if len(opts.Rules) != 0 {
		t.Fatal("expected no rules by default")
	}
}

func TestBuilder_OnMessage_AddsRule(t *testing.T) {
	opts := label.NewBuilder().OnMessage(`error`, "err").Build()
	if len(opts.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(opts.Rules))
	}
	if opts.Rules[0].Label != "err" {
		t.Fatalf("expected label 'err', got %q", opts.Rules[0].Label)
	}
	if opts.Rules[0].Field != "" {
		t.Fatal("expected empty Field for message rule")
	}
}

func TestBuilder_OnField_AddsRule(t *testing.T) {
	opts := label.NewBuilder().OnField("svc", `^db`, "database").Build()
	if len(opts.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(opts.Rules))
	}
	if opts.Rules[0].Field != "svc" {
		t.Fatalf("expected Field 'svc', got %q", opts.Rules[0].Field)
	}
	if opts.Rules[0].Label != "database" {
		t.Fatalf("expected label 'database', got %q", opts.Rules[0].Label)
	}
}

func TestBuilder_WithLabelField(t *testing.T) {
	opts := label.NewBuilder().WithLabelField("tag").Build()
	if opts.LabelField != "tag" {
		t.Fatalf("expected 'tag', got %q", opts.LabelField)
	}
}

func TestBuilder_WithMulti(t *testing.T) {
	opts := label.NewBuilder().WithMulti().Build()
	if !opts.Multi {
		t.Fatal("expected Multi to be true")
	}
}

func TestBuilder_Chaining_MultipleRules(t *testing.T) {
	opts := label.NewBuilder().
		OnMessage(`timeout`, "slow").
		OnMessage(`panic`, "crash").
		OnField("env", `prod`, "production").
		Build()
	if len(opts.Rules) != 3 {
		t.Fatalf("expected 3 rules, got %d", len(opts.Rules))
	}
}
