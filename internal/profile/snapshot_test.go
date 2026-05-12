package profile

import (
	"testing"
	"time"
)

func makeSnapshotProfile(name string, vars map[string]string) Profile {
	p := Profile{Name: name}
	for k, v := range vars {
		p.Vars = append(p.Vars, Var{Name: k, Value: v})
	}
	return p
}

func TestSnapshotProfile_CapturesVars(t *testing.T) {
	p := makeSnapshotProfile("web", map[string]string{"PORT": "8080", "HOST": "localhost"})
	envMap := map[string]string{"PORT": "8080", "HOST": "localhost"}

	before := time.Now().UTC()
	snap := SnapshotProfile(p, envMap)
	after := time.Now().UTC()

	if snap.ProfileName != "web" {
		t.Errorf("expected profile name 'web', got %q", snap.ProfileName)
	}
	if snap.Vars["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", snap.Vars["PORT"])
	}
	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not in expected range", snap.CapturedAt)
	}
}

func TestSnapshotProfile_IsDeepCopy(t *testing.T) {
	p := makeSnapshotProfile("db", map[string]string{"DSN": "postgres://localhost"})
	envMap := map[string]string{"DSN": "postgres://localhost"}

	snap := SnapshotProfile(p, envMap)
	envMap["DSN"] = "mutated"

	if snap.Vars["DSN"] != "postgres://localhost" {
		t.Error("snapshot was mutated when original env map changed")
	}
}

func TestDiffSnapshot_Added(t *testing.T) {
	snap := Snapshot{
		ProfileName: "svc",
		Vars:        map[string]string{"A": "1"},
	}
	current := map[string]string{"A": "1", "B": "2"}

	added, removed, changed := DiffSnapshot(snap, current)

	if len(added) != 1 || added[0] != "B" {
		t.Errorf("expected added=[B], got %v", added)
	}
	if len(removed) != 0 {
		t.Errorf("expected no removed, got %v", removed)
	}
	if len(changed) != 0 {
		t.Errorf("expected no changed, got %v", changed)
	}
}

func TestDiffSnapshot_Removed(t *testing.T) {
	snap := Snapshot{Vars: map[string]string{"A": "1", "B": "2"}}
	current := map[string]string{"A": "1"}

	_, removed, _ := DiffSnapshot(snap, current)
	if len(removed) != 1 || removed[0] != "B" {
		t.Errorf("expected removed=[B], got %v", removed)
	}
}

func TestDiffSnapshot_Changed(t *testing.T) {
	snap := Snapshot{Vars: map[string]string{"A": "old"}}
	current := map[string]string{"A": "new"}

	_, _, changed := DiffSnapshot(snap, current)
	if len(changed) != 1 || changed[0] != "A" {
		t.Errorf("expected changed=[A], got %v", changed)
	}
}

func TestSnapshot_String(t *testing.T) {
	snap := Snapshot{
		ProfileName: "myprofile",
		CapturedAt:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		Vars:        map[string]string{"X": "1", "Y": "2"},
	}
	s := snap.String()
	if s == "" {
		t.Error("String() returned empty")
	}
}
