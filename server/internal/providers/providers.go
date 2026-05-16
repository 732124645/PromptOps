// Package providers is a small model gateway: it dispatches a rendered prompt
// to an LLM provider and returns the completion.
package providers

import "fmt"

// Request is a single completion request.
type Request struct {
	Model   string
	APIKey  string
	BaseURL string
	Prompt  string
}

// Result is a completion response.
type Result struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
	Output   string `json:"output"`
}

// Provider runs a prompt against a specific LLM backend.
type Provider interface {
	Name() string
	Run(req Request) (Result, error)
}

// names is the registry's display order.
var names = []string{"mock", "openai", "claude", "ollama", "gemini"}

var registry = map[string]Provider{
	"mock":   mockProvider{},
	"openai": remoteProvider{kind: "openai"},
	"claude": remoteProvider{kind: "claude"},
	"ollama": remoteProvider{kind: "ollama"},
	"gemini": remoteProvider{kind: "gemini"},
}

// Get returns the named provider, defaulting to "mock" when name is empty.
func Get(name string) (Provider, error) {
	if name == "" {
		name = "mock"
	}
	p, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider %q", name)
	}
	return p, nil
}

// Names returns the registered provider names in display order.
func Names() []string {
	return append([]string(nil), names...)
}

func def(v, fallback string) string {
	if v == "" {
		return fallback
	}
	return v
}
