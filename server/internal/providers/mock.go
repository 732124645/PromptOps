package providers

import (
	"fmt"
	"strings"
)

// mockProvider returns a deterministic completion without any network call.
// It lets the Playground work offline and keeps CI hermetic.
type mockProvider struct{}

func (mockProvider) Name() string { return "mock" }

func (mockProvider) Run(req Request) (Result, error) {
	model := def(req.Model, "mock-1")
	lines := strings.Split(strings.TrimSpace(req.Prompt), "\n")
	first := ""
	if len(lines) > 0 {
		first = lines[0]
	}
	output := fmt.Sprintf(
		"[mock completion]\nmodel: %s\nprompt: %d chars, %d line(s)\nfirst line: %s\n\n"+
			"This is a deterministic mock response. Select a real provider and "+
			"supply an API key to call an actual LLM.",
		model, len(req.Prompt), len(lines), first,
	)
	return Result{Provider: "mock", Model: model, Output: output}, nil
}
