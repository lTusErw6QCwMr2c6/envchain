package exec

import (
	"testing"

	"github.com/envchain/envchain/internal/profile"
	"github.com/envchain/envchain/internal/secret"
)

func setupChainStore(t *testing.T) *profile.Store {
	t.Helper()
	store := newTempStore(t)

	p1 := &profile.Profile{
		Name: "base",
		Vars: []profile.Var{{Name: "APP_ENV"}, {Name: "LOG_LEVEL"}},
	}
	p2 := &profile.Profile{
		Name: "override",
		Vars: []profile.Var{{Name: "LOG_LEVEL"}, {Name: "DEBUG"}},
	}

	if err := store.Save(p1); err != nil {
		t.Fatalf("save base profile: %v", err)
	}
	if err := store.Save(p2); err != nil {
		t.Fatalf("save override profile: %v", err)
	}
	return store
}

func TestChainRunner_ResolveChain_MergesProfiles(t *testing.T) {
	store := setupChainStore(t)
	prov := secret.NewEnvProvider()

	_ = prov.Set("base", "APP_ENV", "production")
	_ = prov.Set("base", "LOG_LEVEL", "info")
	_ = prov.Set("override", "LOG_LEVEL", "debug")
	_ = prov.Set("override", "DEBUG", "true")

	cr := NewChainRunner(store, prov)
	env, err := cr.ResolveChain([]string{"base", "override"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if env["APP_ENV"] != "production" {
		t.Errorf("APP_ENV: got %q, want %q", env["APP_ENV"], "production")
	}
	if env["LOG_LEVEL"] != "debug" {
		t.Errorf("LOG_LEVEL: got %q, want %q (override should win)", env["LOG_LEVEL"], "debug")
	}
	if env["DEBUG"] != "true" {
		t.Errorf("DEBUG: got %q, want %q", env["DEBUG"], "true")
	}
}

func TestChainRunner_ResolveChain_ProfileNotFound(t *testing.T) {
	store := newTempStore(t)
	prov := secret.NewEnvProvider()
	cr := NewChainRunner(store, prov)

	_, err := cr.ResolveChain([]string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for missing profile, got nil")
	}
}

func TestChainRunner_ResolveChain_Empty(t *testing.T) {
	store := newTempStore(t)
	prov := secret.NewEnvProvider()
	cr := NewChainRunner(store, prov)

	_, err := cr.ResolveChain([]string{})
	if err == nil {
		t.Fatal("expected error for empty profile list, got nil")
	}
}
