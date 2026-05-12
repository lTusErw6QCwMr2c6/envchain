package profile_test

import (
	"testing"
	"time"

	"github.com/nicholasgasior/envchain/internal/profile"
)

func newLockStore(t *testing.T) (*profile.LockStore, profile.Store) {
	t.Helper()
	s := newTempStore(t)
	return profile.NewLockStore(s), s
}

func TestLockStore_AcquireAndRelease(t *testing.T) {
	ls, _ := newLockStore(t)

	if err := ls.Acquire("myprofile", "alice", 5*time.Minute); err != nil {
		t.Fatalf("expected acquire to succeed, got: %v", err)
	}

	entry, err := ls.Get("myprofile")
	if err != nil {
		t.Fatalf("expected to get lock, got: %v", err)
	}
	if entry.Owner != "alice" {
		t.Errorf("expected owner alice, got %q", entry.Owner)
	}
	if entry.IsExpired() {
		t.Error("expected lock to not be expired")
	}

	if err := ls.Release("myprofile", "alice"); err != nil {
		t.Fatalf("expected release to succeed, got: %v", err)
	}

	if _, err := ls.Get("myprofile"); err == nil {
		t.Error("expected lock to be gone after release")
	}
}

func TestLockStore_Acquire_BlocksSecondOwner(t *testing.T) {
	ls, _ := newLockStore(t)

	if err := ls.Acquire("myprofile", "alice", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := ls.Acquire("myprofile", "bob", 5*time.Minute)
	if err == nil {
		t.Fatal("expected error when acquiring already-locked profile")
	}
}

func TestLockStore_Acquire_AllowsAfterExpiry(t *testing.T) {
	ls, _ := newLockStore(t)

	if err := ls.Acquire("myprofile", "alice", -1*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Lock is already expired; bob should be able to acquire
	if err := ls.Acquire("myprofile", "bob", 5*time.Minute); err != nil {
		t.Fatalf("expected acquire after expiry to succeed, got: %v", err)
	}
}

func TestLockStore_Release_WrongOwner(t *testing.T) {
	ls, _ := newLockStore(t)

	if err := ls.Acquire("myprofile", "alice", 5*time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err := ls.Release("myprofile", "bob")
	if err == nil {
		t.Fatal("expected error when releasing lock owned by someone else")
	}
}

func TestLockStore_Get_NotFound(t *testing.T) {
	ls, _ := newLockStore(t)

	_, err := ls.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing lock")
	}
}

func TestLockEntry_IsExpired(t *testing.T) {
	expired := profile.LockEntry{ExpiresAt: time.Now().Add(-1 * time.Minute)}
	if !expired.IsExpired() {
		t.Error("expected lock to be expired")
	}

	active := profile.LockEntry{ExpiresAt: time.Now().Add(5 * time.Minute)}
	if active.IsExpired() {
		t.Error("expected lock to not be expired")
	}
}
