package profile_test

import (
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newSearchStore(t *testing.T) *profile.Store {
	t.Helper()
	s, err := profile.NewStore(t.TempDir())
	if err != nil {
		t.Fatalf("newSearchStore: %v", err)
	}
	return s
}

func saveSearchProfile(t *testing.T, s *profile.Store, p profile.Profile) {
	t.Helper()
	if err := s.Save(p); err != nil {
		t.Fatalf("saveSearchProfile: %v", err)
	}
}

func TestSearchProfiles_MatchByName(t *testing.T) {
	s := newSearchStore(t)
	saveSearchProfile(t, s, profile.Profile{Name: "frontend-prod", Vars: []profile.Var{{Key: "PORT", Value: "3000"}}})
	saveSearchProfile(t, s, profile.Profile{Name: "backend-prod", Vars: []profile.Var{{Key: "DB_URL", Value: "postgres://"}}})
	saveSearchProfile(t, s, profile.Profile{Name: "staging", Vars: []profile.Var{{Key: "DEBUG", Value: "true"}}})

	results, err := profile.SearchProfiles(s, profile.SearchOptions{Query: "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.MatchedOn != "name" {
			t.Errorf("expected MatchedOn=name, got %q", r.MatchedOn)
		}
	}
}

func TestSearchProfiles_MatchByVarKey(t *testing.T) {
	s := newSearchStore(t)
	saveSearchProfile(t, s, profile.Profile{Name: "alpha", Vars: []profile.Var{{Key: "DATABASE_URL", Value: "sqlite"}}})
	saveSearchProfile(t, s, profile.Profile{Name: "beta", Vars: []profile.Var{{Key: "PORT", Value: "8080"}}})

	results, err := profile.SearchProfiles(s, profile.SearchOptions{Query: "database"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].MatchedOn != "var_key" {
		t.Errorf("expected var_key, got %q", results[0].MatchedOn)
	}
}

func TestSearchProfiles_MatchByVarValue(t *testing.T) {
	s := newSearchStore(t)
	saveSearchProfile(t, s, profile.Profile{Name: "svc", Vars: []profile.Var{{Key: "ENDPOINT", Value: "https://api.example.com"}}})

	results, err := profile.SearchProfiles(s, profile.SearchOptions{Query: "example.com"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].MatchedOn != "var_value" {
		t.Errorf("expected var_value match, got %+v", results)
	}
}

func TestSearchProfiles_EmptyQuery_ReturnsAll(t *testing.T) {
	s := newSearchStore(t)
	saveSearchProfile(t, s, profile.Profile{Name: "p1", Vars: []profile.Var{{Key: "A", Value: "1"}}})
	saveSearchProfile(t, s, profile.Profile{Name: "p2", Vars: []profile.Var{{Key: "B", Value: "2"}}})

	results, err := profile.SearchProfiles(s, profile.SearchOptions{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2, got %d", len(results))
	}
}

func TestSearchProfiles_CaseSensitive_NoMatch(t *testing.T) {
	s := newSearchStore(t)
	saveSearchProfile(t, s, profile.Profile{Name: "MyProfile", Vars: []profile.Var{{Key: "KEY", Value: "val"}}})

	results, err := profile.SearchProfiles(s, profile.SearchOptions{Query: "myprofile", CaseSensitive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected no results for case-sensitive mismatch, got %d", len(results))
	}
}
