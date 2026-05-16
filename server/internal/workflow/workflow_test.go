package workflow

import (
	"strings"
	"testing"
)

func TestRunChainsStepOutputs(t *testing.T) {
	steps := []Step{
		{Name: "build", Type: "render", Template: "Hello {{name}}"},
		{Name: "model", Type: "model", Provider: "mock", Model: "mock-1"},
		{Name: "upper", Type: "transform", Op: "upper"},
	}
	res, err := Run(steps, map[string]string{"name": "Ada"}, "")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if len(res.Steps) != 3 {
		t.Fatalf("want 3 step results, got %d", len(res.Steps))
	}
	if res.Steps[0].Output != "Hello Ada" {
		t.Fatalf("render step = %q", res.Steps[0].Output)
	}
	if !strings.Contains(res.Steps[1].Output, "mock completion") {
		t.Fatalf("model step = %q", res.Steps[1].Output)
	}
	if res.Output != strings.ToUpper(res.Output) {
		t.Fatalf("final output not upper-cased: %q", res.Output)
	}
}

func TestRunRenderUsesPreviousOutput(t *testing.T) {
	steps := []Step{
		{Name: "first", Type: "render", Template: "value"},
		{Name: "second", Type: "render", Template: "[{{input}}]"},
	}
	res, err := Run(steps, nil, "")
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if res.Output != "[value]" {
		t.Fatalf("output = %q, want [value]", res.Output)
	}
}

func TestRunRejectsUnknownStepType(t *testing.T) {
	_, err := Run([]Step{{Name: "x", Type: "bogus"}}, nil, "")
	if err == nil {
		t.Fatal("expected error for unknown step type")
	}
}
