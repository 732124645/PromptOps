// Package workflow runs ordered prompt/model/transform pipelines.
package workflow

import (
	"fmt"
	"strings"

	"github.com/732124645/promptops/server/internal/providers"
	"github.com/732124645/promptops/server/internal/render"
)

// Step is one stage of a workflow.
//
// Types:
//   - "render":    render Template (with run variables plus {{input}}, the
//     previous step's output) into the new output.
//   - "model":     send the current input to Provider/Model and take the
//     completion as the new output.
//   - "transform": apply Op ("upper", "lower", "trim") to the current input.
type Step struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Template string `json:"template,omitempty"`
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
	Op       string `json:"op,omitempty"`
}

// StepResult is the recorded output of one executed step.
type StepResult struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Output string `json:"output"`
}

// Result is the outcome of a workflow run.
type Result struct {
	Output string       `json:"output"`
	Steps  []StepResult `json:"steps"`
}

// Run executes steps in order, threading each step's output into the next.
func Run(steps []Step, vars map[string]string, apiKey string) (Result, error) {
	input := ""
	trace := make([]StepResult, 0, len(steps))

	for i, step := range steps {
		var output string
		switch step.Type {
		case "render":
			ctx := map[string]string{"input": input}
			for k, v := range vars {
				ctx[k] = v
			}
			output = render.Template(step.Template, ctx)
		case "model":
			prov, err := providers.Get(step.Provider)
			if err != nil {
				return Result{}, fmt.Errorf("step %d (%s): %w", i+1, step.Name, err)
			}
			res, err := prov.Run(providers.Request{
				Model:  step.Model,
				APIKey: apiKey,
				Prompt: input,
			})
			if err != nil {
				return Result{}, fmt.Errorf("step %d (%s): %w", i+1, step.Name, err)
			}
			output = res.Output
		case "transform":
			output = transform(step.Op, input)
		default:
			return Result{}, fmt.Errorf("step %d (%s): unknown type %q", i+1, step.Name, step.Type)
		}
		input = output
		trace = append(trace, StepResult{Name: step.Name, Type: step.Type, Output: output})
	}

	return Result{Output: input, Steps: trace}, nil
}

func transform(op, s string) string {
	switch op {
	case "upper":
		return strings.ToUpper(s)
	case "lower":
		return strings.ToLower(s)
	case "trim":
		return strings.TrimSpace(s)
	default:
		return s
	}
}
