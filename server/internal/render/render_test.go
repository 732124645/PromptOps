package render

import (
	"reflect"
	"testing"
)

func TestTemplate(t *testing.T) {
	cases := []struct {
		name    string
		content string
		vars    map[string]string
		want    string
	}{
		{"substitutes", "hi {{name}}", map[string]string{"name": "Ada"}, "hi Ada"},
		{"keeps unknown", "hi {{name}}", map[string]string{}, "hi {{name}}"},
		{"whitespace and dotted", "{{ a.b }}", map[string]string{"a.b": "X"}, "X"},
		{"repeated", "{{x}}-{{x}}", map[string]string{"x": "1"}, "1-1"},
	}
	for _, tc := range cases {
		if got := Template(tc.content, tc.vars); got != tc.want {
			t.Errorf("%s: Template() = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestVariables(t *testing.T) {
	got := Variables("you are a {{role}} expert, {{name}}. {{role}} again.")
	want := []string{"role", "name"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Variables() = %v, want %v", got, want)
	}
}
