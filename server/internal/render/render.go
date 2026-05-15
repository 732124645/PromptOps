package render

import "regexp"

var varRE = regexp.MustCompile(`{{\s*([\w.]+)\s*}}`)

// Template substitutes {{variable}} placeholders with values from vars.
// Unknown variables are left untouched.
func Template(content string, vars map[string]string) string {
	return varRE.ReplaceAllStringFunc(content, func(match string) string {
		key := varRE.FindStringSubmatch(match)[1]
		if v, ok := vars[key]; ok {
			return v
		}
		return match
	})
}

// Variables returns the distinct variable names referenced in content.
func Variables(content string) []string {
	seen := map[string]bool{}
	var out []string
	for _, m := range varRE.FindAllStringSubmatch(content, -1) {
		if !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}
