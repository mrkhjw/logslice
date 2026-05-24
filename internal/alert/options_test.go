package alert_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/alert"
)

func TestBuilder_Defaults_Empty(t *testing.T) {
	rules := alert.NewBuilder().Build()
	if len(rules) != 0 {
		t.Errorf("expected 0 rules, got %d", len(rules))
	}
}

func TestBuilder_OnLevel_AddsRule(t *testing.T) {
	rules := alert.NewBuilder().
		OnLevel("warn", 10, time.Minute).
		Build()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Level != "warn" {
		t.Errorf("expected level warn, got %s", rules[0].Level)
	}
	if rules[0].Threshold != 10 {
		t.Errorf("expected threshold 10, got %d", rules[0].Threshold)
	}
	if rules[0].Window != time.Minute {
		t.Errorf("expected window 1m, got %s", rules[0].Window)
	}
}

func TestBuilder_OnError_SetsLevel(t *testing.T) {
	rules := alert.NewBuilder().OnError(5, 30*time.Second).Build()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Level != "error" {
		t.Errorf("expected level error, got %s", rules[0].Level)
	}
}

func TestBuilder_OnFatal_SetsLevel(t *testing.T) {
	rules := alert.NewBuilder().OnFatal(1, time.Hour).Build()
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	if rules[0].Level != "fatal" {
		t.Errorf("expected level fatal, got %s", rules[0].Level)
	}
}

func TestBuilder_MultipleRules(t *testing.T) {
	rules := alert.NewBuilder().
		OnError(3, time.Minute).
		OnFatal(1, time.Minute).
		OnLevel("warn", 20, 5*time.Minute).
		Build()
	if len(rules) != 3 {
		t.Errorf("expected 3 rules, got %d", len(rules))
	}
}

func TestBuilder_Build_ReturnsCopy(t *testing.T) {
	b := alert.NewBuilder().OnError(1, time.Minute)
	r1 := b.Build()
	r2 := b.Build()
	if &r1[0] == &r2[0] {
		t.Error("Build should return independent copies")
	}
}
