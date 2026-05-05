package profile

import (
	"strings"
	"testing"
)

func TestExpandVars_WithEnvMap(t *testing.T) {
	p := &Profile{
		Name: "test",
		Vars: []Var{
			{Key: "URL", Value: "http://${HOST}:${PORT}"},
			{Key: "PLAIN", Value: "no-expand"},
		},
	}
	env := map[string]string{
		"HOST": "localhost",
		"PORT": "8080",
	}
	expanded := ExpandVars(p, env)
	if expanded.Vars[0].Value != "http://localhost:8080" {
		t.Errorf("expected expanded URL, got %q", expanded.Vars[0].Value)
	}
	if expanded.Vars[1].Value != "no-expand" {
		t.Errorf("expected plain value unchanged, got %q", expanded.Vars[1].Value)
	}
}

func TestExpandVars_DoesNotMutateOriginal(t *testing.T) {
	p := &Profile{
		Name: "test",
		Vars: []Var{{Key: "X", Value: "${Y}"}},
	}
	env := map[string]string{"Y": "replaced"}
	ExpandVars(p, env)
	if p.Vars[0].Value != "${Y}" {
		t.Error("original profile was mutated")
	}
}

func TestRenderTemplate_Dotenv(t *testing.T) {
	p := &Profile{
		Name: "test",
		Vars: []Var{
			{Key: "FOO", Value: "bar"},
			{Key: "BAZ", Value: "qux"},
		},
	}
	out, err := RenderTemplate(p, "dotenv")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "FOO=bar\n") {
		t.Errorf("expected FOO=bar in output, got: %q", out)
	}
	if !strings.Contains(out, "BAZ=qux\n") {
		t.Errorf("expected BAZ=qux in output, got: %q", out)
	}
}

func TestRenderTemplate_Export(t *testing.T) {
	p := &Profile{
		Name: "test",
		Vars: []Var{{Key: "FOO", Value: "bar baz"}},
	}
	out, err := RenderTemplate(p, "export")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "export FOO=") {
		t.Errorf("expected export statement, got: %q", out)
	}
}

func TestRenderTemplate_UnknownFormat(t *testing.T) {
	p := &Profile{Name: "test"}
	_, err := RenderTemplate(p, "xml")
	if err == nil {
		t.Error("expected error for unknown format")
	}
}
