package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newCloneStore(t *testing.T) *profile.Store {
	t.Helper()
	st, err := profile.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func saveCloneProfile(t *testing.T, st *profile.Store, p *profile.Profile) {
	t.Helper()
	if err := st.Save(p); err != nil {
		t.Fatalf("Save(%q): %v", p.Name, err)
	}
}

func TestCloneProfile_Basic(t *testing.T) {
	st := newCloneStore(t)
	saveCloneProfile(t, st, &profile.Profile{
		Name: "src",
		Vars: []profile.Var{{Key: "FOO", Value: "bar"}, {Key: "BAZ", Value: "qux"}},
	})

	if err := profile.CloneProfile(st, "src", "dst", profile.CloneOptions{}); err != nil {
		t.Fatalf("CloneProfile: %v", err)
	}

	dst, err := st.Load("dst")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	if len(dst.Vars) != 2 {
		t.Errorf("expected 2 vars, got %d", len(dst.Vars))
	}
}

func TestCloneProfile_PreservesParents(t *testing.T) {
	st := newCloneStore(t)
	saveCloneProfile(t, st, &profile.Profile{Name: "base"})
	saveCloneProfile(t, st, &profile.Profile{
		Name:    "src",
		Parents: []string{"base"},
		Vars:    []profile.Var{{Key: "X", Value: "1"}},
	})

	if err := profile.CloneProfile(st, "src", "dst", profile.CloneOptions{}); err != nil {
		t.Fatalf("CloneProfile: %v", err)
	}

	dst, err := st.Load("dst")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	if len(dst.Parents) != 1 || dst.Parents[0] != "base" {
		t.Errorf("expected parents [base], got %v", dst.Parents)
	}
}

func TestCloneProfile_StripParents(t *testing.T) {
	st := newCloneStore(t)
	saveCloneProfile(t, st, &profile.Profile{Name: "base"})
	saveCloneProfile(t, st, &profile.Profile{
		Name:    "src",
		Parents: []string{"base"},
		Vars:    []profile.Var{{Key: "X", Value: "1"}},
	})

	if err := profile.CloneProfile(st, "src", "dst", profile.CloneOptions{StripParents: true}); err != nil {
		t.Fatalf("CloneProfile: %v", err)
	}

	dst, err := st.Load("dst")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	if len(dst.Parents) != 0 {
		t.Errorf("expected no parents after strip, got %v", dst.Parents)
	}
}

func TestCloneProfile_NoDuplicateWithoutOverwrite(t *testing.T) {
	st := newCloneStore(t)
	saveCloneProfile(t, st, &profile.Profile{Name: "src", Vars: []profile.Var{{Key: "A", Value: "1"}}})
	saveCloneProfile(t, st, &profile.Profile{Name: "dst", Vars: []profile.Var{{Key: "B", Value: "2"}}})

	err := profile.CloneProfile(st, "src", "dst", profile.CloneOptions{})
	if err == nil {
		t.Fatal("expected error when dst exists without overwrite, got nil")
	}
}

func TestCloneProfile_OverwriteAllowed(t *testing.T) {
	st := newCloneStore(t)
	saveCloneProfile(t, st, &profile.Profile{Name: "src", Vars: []profile.Var{{Key: "NEW", Value: "val"}}})
	saveCloneProfile(t, st, &profile.Profile{Name: "dst", Vars: []profile.Var{{Key: "OLD", Value: "old"}}})

	if err := profile.CloneProfile(st, "src", "dst", profile.CloneOptions{Overwrite: true}); err != nil {
		t.Fatalf("CloneProfile with overwrite: %v", err)
	}

	dst, err := st.Load("dst")
	if err != nil {
		t.Fatalf("Load dst: %v", err)
	}
	if len(dst.Vars) != 1 || dst.Vars[0].Key != "NEW" {
		t.Errorf("expected var NEW after overwrite, got %v", dst.Vars)
	}
}
