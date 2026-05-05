package profile

import (
	"fmt"
	"os"
	"strings"
)

// ExpandVars expands environment variable references within profile variable
// values using the provided env map as the source of substitutions.
// References use the ${VAR} or $VAR syntax.
func ExpandVars(p *Profile, env map[string]string) *Profile {
	expanded := &Profile{
		Name:    p.Name,
		Parents: p.Parents,
		Vars:    make([]Var, len(p.Vars)),
	}

	mapper := func(key string) string {
		if val, ok := env[key]; ok {
			return val
		}
		return os.Getenv(key)
	}

	for i, v := range p.Vars {
		expanded.Vars[i] = Var{
			Key:   v.Key,
			Value: os.Expand(v.Value, mapper),
		}
	}
	return expanded
}

// RenderTemplate renders a simple key=value template for a profile,
// suitable for use in shell scripts or .env files.
func RenderTemplate(p *Profile, format string) (string, error) {
	switch strings.ToLower(format) {
	case "dotenv", "":
		return renderDotenv(p), nil
	case "export":
		return renderExport(p), nil
	default:
		return "", fmt.Errorf("unknown template format: %q", format)
	}
}

func renderDotenv(p *Profile) string {
	var sb strings.Builder
	for _, v := range p.Vars {
		fmt.Fprintf(&sb, "%s=%s\n", v.Key, v.Value)
	}
	return sb.String()
}

func renderExport(p *Profile) string {
	var sb strings.Builder
	for _, v := range p.Vars {
		fmt.Fprintf(&sb, "export %s=%q\n", v.Key, v.Value)
	}
	return sb.String()
}
