package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newRenameStore(t *testing.T) *profile.Store {
	t.Helper()
	st, err := profile.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("newRenameStore: %v", err)
	}
	return st
}

func saveRenameProfile(t *testing.T, st *profile.Store, name string) {
	t.Helper()
	p := &profile.Profile{
		Name: name,
		Vars: []profile.Var{{Key: "FOO", Value: "bar"}},
	}
	if err := st.Save(p); err != nil {
		t.Fatalf("saveRenameProfile %q: %v", name, err)
	}
}

func TestRenameProfile_Basic(t *testing.T) {
	st := newRenameStore(t)
	saveRenameProfile(t, st, "old")

	if err := profile.RenameProfile(st, "old", "new", profile.RenameOptions{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := st.Load("new"); err != nil {
		t.Errorf("expected 'new' profile to exist: %v", err)
	}
	if _, err := st.Load("old"); err == nil {
		t.Errorf("expected 'old' profile to be deleted")
	}
}

func TestRenameProfile_DestinationExists_NoOverwrite(t *testing.T) {
	st := newRenameStore(t)
	saveRenameProfile(t, st, "alpha")
	saveRenameProfile(t, st, "beta")

	err := profile.RenameProfile(st, "alpha", "beta", profile.RenameOptions{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when destination exists without overwrite flag")
	}
}

func TestRenameProfile_DestinationExists_WithOverwrite(t *testing.T) {
	st := newRenameStore(t)
	saveRenameProfile(t, st, "alpha")
	saveRenameProfile(t, st, "beta")

	if err := profile.RenameProfile(st, "alpha", "beta", profile.RenameOptions{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	p, err := st.Load("beta")
	if err != nil {
		t.Fatalf("load beta: %v", err)
	}
	if p.Name != "beta" {
		t.Errorf("expected name 'beta', got %q", p.Name)
	}
}

func TestRenameProfile_SourceNotFound(t *testing.T) {
	st := newRenameStore(t)

	err := profile.RenameProfile(st, "ghost", "target", profile.RenameOptions{})
	if err == nil {
		t.Fatal("expected error for missing source profile")
	}
}

func TestRenameProfile_InvalidDestinationName(t *testing.T) {
	st := newRenameStore(t)
	saveRenameProfile(t, st, "src")

	err := profile.RenameProfile(st, "src", "bad name!", profile.RenameOptions{})
	if err == nil {
		t.Fatal("expected error for invalid destination name")
	}
}
