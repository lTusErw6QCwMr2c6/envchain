package profile

import (
	"testing"
)

func TestLintProfile_NoVars(t *testing.T) {
	p := &Profile{Name: "empty"}
	res := LintProfile(p)
	if res.OK() {
		t.Error("expected warning for empty vars")
	}
	if len(res.Warnings) != 1 {
		t.Errorf("expected 1 warning, got %d", len(res.Warnings))
	}
}

func TestLintProfile_LowercaseVar(t *testing.T) {
	p := &Profile{
		Name: "myprofile",
		Vars: []Var{{Name: "foo"}, {Name: "BAR"}},
	}
	res := LintProfile(p)
	if res.OK() {
		t.Error("expected warning for lowercase variable")
	}
	found := false
	for _, w := range res.Warnings {
		if w.Field == "vars.foo" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning about 'foo' variable")
	}
}

func TestLintProfile_UnderscorePrefix(t *testing.T) {
	p := &Profile{
		Name: "myprofile",
		Vars: []Var{{Name: "_INTERNAL"}},
	}
	res := LintProfile(p)
	if res.OK() {
		t.Error("expected warning for underscore-prefixed variable")
	}
}

func TestLintProfile_TooManyParents(t *testing.T) {
	p := &Profile{
		Name:    "deep",
		Vars:    []Var{{Name: "FOO"}},
		Parents: []string{"a", "b", "c", "d", "e", "f"},
	}
	res := LintProfile(p)
	if res.OK() {
		t.Error("expected warning for too many parents")
	}
}

func TestLintProfile_Clean(t *testing.T) {
	p := &Profile{
		Name: "clean",
		Vars: []Var{{Name: "FOO"}, {Name: "BAR"}},
	}
	res := LintProfile(p)
	if !res.OK() {
		t.Errorf("expected no warnings, got: %v", res.Warnings)
	}
}

func TestLintWarning_String(t *testing.T) {
	w := LintWarning{Field: "vars.foo", Message: "not uppercase"}
	s := w.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
