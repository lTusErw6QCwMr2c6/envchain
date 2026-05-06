package profile

import (
	"testing"
)

func makeProfile(vars []EnvVar) *Profile {
	return &Profile{Name: "test", Vars: vars}
}

func TestDiffProfiles_NoDiff(t *testing.T) {
	base := makeProfile([]EnvVar{{Key: "FOO", Value: "bar"}})
	target := makeProfile([]EnvVar{{Key: "FOO", Value: "bar"}})

	entries := DiffProfiles(base, target)
	if len(entries) != 0 {
		t.Fatalf("expected no diff, got %d entries", len(entries))
	}
}

func TestDiffProfiles_Added(t *testing.T) {
	base := makeProfile([]EnvVar{})
	target := makeProfile([]EnvVar{{Key: "NEW_VAR", Value: "hello"}})

	entries := DiffProfiles(base, target)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Kind != DiffAdded || entries[0].Key != "NEW_VAR" || entries[0].New != "hello" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiffProfiles_Removed(t *testing.T) {
	base := makeProfile([]EnvVar{{Key: "OLD_VAR", Value: "bye"}})
	target := makeProfile([]EnvVar{})

	entries := DiffProfiles(base, target)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Kind != DiffRemoved || entries[0].Key != "OLD_VAR" || entries[0].Old != "bye" {
		t.Errorf("unexpected entry: %+v", entries[0])
	}
}

func TestDiffProfiles_Changed(t *testing.T) {
	base := makeProfile([]EnvVar{{Key: "HOST", Value: "localhost"}})
	target := makeProfile([]EnvVar{{Key: "HOST", Value: "prod.example.com"}})

	entries := DiffProfiles(base, target)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Kind != DiffChanged || e.Key != "HOST" || e.Old != "localhost" || e.New != "prod.example.com" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestDiffEntry_String(t *testing.T) {
	cases := []struct {
		entry    DiffEntry
		expected string
	}{
		{DiffEntry{Key: "A", Kind: DiffAdded, New: "1"}, "+ A=1"},
		{DiffEntry{Key: "B", Kind: DiffRemoved, Old: "2"}, "- B=2"},
		{DiffEntry{Key: "C", Kind: DiffChanged, Old: "x", New: "y"}, "~ C: x -> y"},
	}
	for _, tc := range cases {
		if got := tc.entry.String(); got != tc.expected {
			t.Errorf("String() = %q, want %q", got, tc.expected)
		}
	}
}
