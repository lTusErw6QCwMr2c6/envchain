package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newPromoteStore(t *testing.T) profile.Store {
	t.Helper()
	st, err := profile.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("newPromoteStore: %v", err)
	}
	return st
}

func savePromoteProfile(t *testing.T, st profile.Store, name string, parents []string, vars ...profile.Var) {
	t.Helper()
	p := profile.Profile{Name: name, Parents: parents, Vars: vars}
	if err := st.Save(p); err != nil {
		t.Fatalf("savePromoteProfile %q: %v", name, err)
	}
}

func TestPromoteProfile_Basic(t *testing.T) {
	st := newPromoteStore(t)
	savePromoteProfile(t, st, "staging", nil, profile.Var{Key: "DB_HOST", Value: "staging-db"})

	if err := profile.PromoteProfile(st, "staging", "production", profile.PromoteOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, err := st.Load("production")
	if err != nil {
		t.Fatalf("load production: %v", err)
	}
	if len(dst.Vars) != 1 || dst.Vars[0].Key != "DB_HOST" {
		t.Errorf("expected DB_HOST in production, got %+v", dst.Vars)
	}
}

func TestPromoteProfile_NoOverwrite(t *testing.T) {
	st := newPromoteStore(t)
	savePromoteProfile(t, st, "staging", nil, profile.Var{Key: "X", Value: "1"})
	savePromoteProfile(t, st, "production", nil, profile.Var{Key: "X", Value: "old"})

	err := profile.PromoteProfile(st, "staging", "production", profile.PromoteOptions{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite")
	}
}

func TestPromoteProfile_OverwriteAllowed(t *testing.T) {
	st := newPromoteStore(t)
	savePromoteProfile(t, st, "staging", nil, profile.Var{Key: "X", Value: "new"})
	savePromoteProfile(t, st, "production", nil, profile.Var{Key: "X", Value: "old"})

	if err := profile.PromoteProfile(st, "staging", "production", profile.PromoteOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, err := st.Load("production")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if dst.Vars[0].Value != "new" {
		t.Errorf("expected value 'new', got %q", dst.Vars[0].Value)
	}
}

func TestPromoteProfile_StripParents(t *testing.T) {
	st := newPromoteStore(t)
	savePromoteProfile(t, st, "base", nil, profile.Var{Key: "BASE", Value: "1"})
	savePromoteProfile(t, st, "staging", []string{"base"}, profile.Var{Key: "APP", Value: "staging"})

	if err := profile.PromoteProfile(st, "staging", "production", profile.PromoteOptions{StripParents: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dst, err := st.Load("production")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(dst.Parents) != 0 {
		t.Errorf("expected no parents, got %v", dst.Parents)
	}
}

func TestPromoteProfile_SameSrcDst(t *testing.T) {
	st := newPromoteStore(t)
	savePromoteProfile(t, st, "staging", nil)

	err := profile.PromoteProfile(st, "staging", "staging", profile.PromoteOptions{})
	if err == nil {
		t.Fatal("expected error for same src and dst")
	}
}
