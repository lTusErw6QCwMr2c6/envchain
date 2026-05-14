package profile_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/nicholasgasior/envchain/internal/profile"
)

// inMemoryArchiveStore is a simple in-memory ArchiveStore for tests.
type inMemoryArchiveStore struct {
	mu       sync.RWMutex
	archives map[string]profile.ArchivedProfile
}

func newInMemoryArchiveStore() *inMemoryArchiveStore {
	return &inMemoryArchiveStore{archives: make(map[string]profile.ArchivedProfile)}
}

func (s *inMemoryArchiveStore) SaveArchive(name string, a profile.ArchivedProfile) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.archives[name] = a
	return nil
}

func (s *inMemoryArchiveStore) LoadArchive(name string) (profile.ArchivedProfile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.archives[name]
	if !ok {
		return profile.ArchivedProfile{}, fmt.Errorf("archive not found: %s", name)
	}
	return a, nil
}

func (s *inMemoryArchiveStore) ListArchives() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	names := make([]string, 0, len(s.archives))
	for k := range s.archives {
		names = append(names, k)
	}
	return names, nil
}

func (s *inMemoryArchiveStore) DeleteArchive(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.archives[name]; !ok {
		return fmt.Errorf("archive not found: %s", name)
	}
	delete(s.archives, name)
	return nil
}

func newArchiveActiveStore(t *testing.T) profile.Store {
	t.Helper()
	dir := t.TempDir()
	st, err := profile.NewStore(dir)
	if err != nil {
		t.Fatalf("newArchiveActiveStore: %v", err)
	}
	return st
}

func TestArchiveProfile_Basic(t *testing.T) {
	active := newArchiveActiveStore(t)
	archive := newInMemoryArchiveStore()

	p := profile.Profile{Name: "staging", Vars: []profile.Var{{Key: "DB_URL", Value: "postgres://localhost"}}}
	if err := active.Save(p); err != nil {
		t.Fatalf("save: %v", err)
	}

	if err := profile.ArchiveProfile(active, archive, "staging", "alice", "decommissioned"); err != nil {
		t.Fatalf("ArchiveProfile: %v", err)
	}

	// Should be gone from active store.
	if _, err := active.Load("staging"); err == nil {
		t.Error("expected profile to be removed from active store")
	}

	a, err := archive.LoadArchive("staging")
	if err != nil {
		t.Fatalf("LoadArchive: %v", err)
	}
	if a.ArchivedBy != "alice" {
		t.Errorf("ArchivedBy = %q, want %q", a.ArchivedBy, "alice")
	}
	if a.Reason != "decommissioned" {
		t.Errorf("Reason = %q, want %q", a.Reason, "decommissioned")
	}
	if a.ArchivedAt.IsZero() {
		t.Error("ArchivedAt should not be zero")
	}
}

func TestRestoreProfile_Basic(t *testing.T) {
	active := newArchiveActiveStore(t)
	archive := newInMemoryArchiveStore()

	p := profile.Profile{Name: "prod", Vars: []profile.Var{{Key: "API_KEY", Value: "secret"}}}
	if err := active.Save(p); err != nil {
		t.Fatalf("save: %v", err)
	}
	if err := profile.ArchiveProfile(active, archive, "prod", "bob", "temp removal"); err != nil {
		t.Fatalf("ArchiveProfile: %v", err)
	}

	restored, err := profile.RestoreProfile(archive, active, "prod", false)
	if err != nil {
		t.Fatalf("RestoreProfile: %v", err)
	}
	if restored.Profile.Name != "prod" {
		t.Errorf("restored profile name = %q, want %q", restored.Profile.Name, "prod")
	}

	loaded, err := active.Load("prod")
	if err != nil {
		t.Fatalf("Load after restore: %v", err)
	}
	if len(loaded.Vars) != 1 || loaded.Vars[0].Key != "API_KEY" {
		t.Errorf("unexpected vars after restore: %+v", loaded.Vars)
	}
}

func TestRestoreProfile_NoOverwrite(t *testing.T) {
	active := newArchiveActiveStore(t)
	archive := newInMemoryArchiveStore()

	p := profile.Profile{Name: "dev", Vars: []profile.Var{{Key: "FOO", Value: "bar"}}}
	if err := active.Save(p); err != nil {
		t.Fatalf("save: %v", err)
	}
	// Manually add to archive without removing from active.
	if err := archive.SaveArchive("dev", profile.ArchivedProfile{Profile: p}); err != nil {
		t.Fatalf("SaveArchive: %v", err)
	}

	_, err := profile.RestoreProfile(archive, active, "dev", false)
	if err == nil {
		t.Error("expected error when restoring over existing profile without overwrite")
	}
}
