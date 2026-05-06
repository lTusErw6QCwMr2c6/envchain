package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newCopyStore(t *testing.T) *profile.Store {
	t.Helper()
	st, err := profile.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return st
}

func saveTestProfile(t *testing.T, st *profile.Store, name string, vars []profile.Var, parents []string) {
	t.Helper()
	p := &profile.Profile{Name: name, Vars: vars, Parents: parents}
	if err := st.Save(p); err != nil {
		t.Fatalf("Save %q: %v", name, err)
	}
}

func TestCopyProfile_Basic(t *testing.T) {
	st := newCopyStore(t)
	saveTestProfile(t, st, "src", []profile.Var{{Name: "FOO", Value: "bar"}}, nil)

	cp, err := profile.CopyProfile(st, "src", "dst", profile.CopyOptions{})
	if err != nil {
		t.Fatalf("CopyProfile: %v", err)
	}
	if cp.Name != "dst" {
		t.Errorf("expected name dst, got %q", cp.Name)
	}
	if len(cp.Vars) != 1 || cp.Vars[0].Name != "FOO" {
		t.Errorf("unexpected vars: %v", cp.Vars)
	}
}

func TestCopyProfile_NoDuplicateWithoutOverwrite(t *testing.T) {
	st := newCopyStore(t)
	saveTestProfile(t, st, "src", []profile.Var{{Name: "X", Value: "1"}}, nil)
	saveTestProfile(t, st, "dst", []profile.Var{{Name: "Y", Value: "2"}}, nil)

	_, err := profile.CopyProfile(st, "src", "dst", profile.CopyOptions{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when destination exists and Overwrite=false")
	}
}

func TestCopyProfile_OverwriteAllowed(t *testing.T) {
	st := newCopyStore(t)
	saveTestProfile(t, st, "src", []profile.Var{{Name: "X", Value: "new"}}, nil)
	saveTestProfile(t, st, "dst", []profile.Var{{Name: "X", Value: "old"}}, nil)

	cp, err := profile.CopyProfile(st, "src", "dst", profile.CopyOptions{Overwrite: true})
	if err != nil {
		t.Fatalf("CopyProfile with Overwrite: %v", err)
	}
	if cp.Vars[0].Value != "new" {
		t.Errorf("expected value 'new', got %q", cp.Vars[0].Value)
	}
}

func TestCopyProfile_ResetParents(t *testing.T) {
	st := newCopyStore(t)
	saveTestProfile(t, st, "base", nil, nil)
	saveTestProfile(t, st, "src", nil, []string{"base"})

	cp, err := profile.CopyProfile(st, "src", "dst", profile.CopyOptions{ResetParents: true})
	if err != nil {
		t.Fatalf("CopyProfile ResetParents: %v", err)
	}
	if len(cp.Parents) != 0 {
		t.Errorf("expected no parents after reset, got %v", cp.Parents)
	}
}

func TestRenameProfile(t *testing.T) {
	st := newCopyStore(t)
	saveTestProfile(t, st, "old", []profile.Var{{Name: "K", Value: "v"}}, nil)

	_, err := profile.RenameProfile(st, "old", "new", profile.CopyOptions{})
	if err != nil {
		t.Fatalf("RenameProfile: %v", err)
	}

	if _, err := st.Load("old"); err == nil {
		t.Error("expected old profile to be deleted")
	}
	loaded, err := st.Load("new")
	if err != nil {
		t.Fatalf("Load renamed profile: %v", err)
	}
	if loaded.Vars[0].Value != "v" {
		t.Errorf("unexpected value: %q", loaded.Vars[0].Value)
	}
}
