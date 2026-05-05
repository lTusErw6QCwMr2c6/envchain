package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newChainStore(t *testing.T) *profile.Store {
	t.Helper()
	dir := t.TempDir()
	store, err := profile.NewStore(dir)
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	return store
}

func saveProfile(t *testing.T, store *profile.Store, p *profile.Profile) {
	t.Helper()
	if err := store.Save(p); err != nil {
		t.Fatalf("Save(%q): %v", p.Name, err)
	}
}

func TestChainNames_SingleProfile(t *testing.T) {
	store := newChainStore(t)
	saveProfile(t, store, &profile.Profile{Name: "base", Vars: []profile.Var{{Key: "A", Value: "1"}}})

	names, err := profile.ChainNames(store, "base")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(names) != 1 || names[0] != "base" {
		t.Errorf("expected [base], got %v", names)
	}
}

func TestChainNames_WithParents(t *testing.T) {
	store := newChainStore(t)
	saveProfile(t, store, &profile.Profile{Name: "root", Vars: []profile.Var{{Key: "ROOT", Value: "r"}}})
	saveProfile(t, store, &profile.Profile{Name: "mid", Chain: []string{"root"}, Vars: []profile.Var{{Key: "MID", Value: "m"}}})
	saveProfile(t, store, &profile.Profile{Name: "leaf", Chain: []string{"mid"}, Vars: []profile.Var{{Key: "LEAF", Value: "l"}}})

	names, err := profile.ChainNames(store, "leaf")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"root", "mid", "leaf"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("names[%d]: want %q, got %q", i, n, names[i])
		}
	}
}

func TestChainNames_CycleDetection(t *testing.T) {
	store := newChainStore(t)
	saveProfile(t, store, &profile.Profile{Name: "a", Chain: []string{"b"}})
	saveProfile(t, store, &profile.Profile{Name: "b", Chain: []string{"a"}})

	_, err := profile.ChainNames(store, "a")
	if err == nil {
		t.Fatal("expected cycle error, got nil")
	}
}

func TestResolveChain_MergesInOrder(t *testing.T) {
	store := newChainStore(t)
	saveProfile(t, store, &profile.Profile{Name: "base", Vars: []profile.Var{{Key: "X", Value: "base"}, {Key: "BASE", Value: "1"}}})
	saveProfile(t, store, &profile.Profile{Name: "override", Vars: []profile.Var{{Key: "X", Value: "override"}}})

	env, err := profile.ResolveChain(store, []string{"base", "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["X"] != "override" {
		t.Errorf("X: want %q, got %q", "override", env["X"])
	}
	if env["BASE"] != "1" {
		t.Errorf("BASE: want %q, got %q", "1", env["BASE"])
	}
}

func TestResolveChain_Empty(t *testing.T) {
	store := newChainStore(t)
	env, err := profile.ResolveChain(store, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 0 {
		t.Errorf("expected empty map, got %v", env)
	}
}
