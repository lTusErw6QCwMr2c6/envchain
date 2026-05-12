package profile

import (
	"fmt"
	"time"
)

// LockEntry represents a lock on a profile, preventing concurrent modifications.
type LockEntry struct {
	Profile   string    `json:"profile"`
	Owner     string    `json:"owner"`
	LockedAt  time.Time `json:"locked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired returns true if the lock has passed its expiry time.
func (l LockEntry) IsExpired() bool {
	return time.Now().After(l.ExpiresAt)
}

// LockStore manages profile locks backed by a key-value store.
type LockStore struct {
	store Store
	ns    string
}

// NewLockStore creates a LockStore using the given profile Store.
func NewLockStore(s Store) *LockStore {
	return &LockStore{store: s, ns: "__lock__"}
}

func (ls *LockStore) lockKey(profile string) string {
	return ls.ns + profile
}

// Acquire attempts to lock the given profile for owner with a TTL duration.
// Returns an error if the profile is already locked by someone else.
func (ls *LockStore) Acquire(profile, owner string, ttl time.Duration) error {
	key := ls.lockKey(profile)
	existing, err := ls.store.Load(key)
	if err == nil {
		entry := toLockEntry(existing)
		if !entry.IsExpired() {
			return fmt.Errorf("profile %q is locked by %q until %s", profile, entry.Owner, entry.ExpiresAt.Format(time.RFC3339))
		}
	}
	now := time.Now().UTC()
	lock := Profile{
		Name: key,
		Vars: map[string]string{
			"OWNER":      owner,
			"LOCKED_AT":  now.Format(time.RFC3339),
			"EXPIRES_AT": now.Add(ttl).Format(time.RFC3339),
		},
	}
	return ls.store.Save(lock)
}

// Release removes the lock for the given profile if owned by owner.
func (ls *LockStore) Release(profile, owner string) error {
	key := ls.lockKey(profile)
	existing, err := ls.store.Load(key)
	if err != nil {
		return fmt.Errorf("no lock found for profile %q", profile)
	}
	entry := toLockEntry(existing)
	if entry.Owner != owner {
		return fmt.Errorf("profile %q is locked by %q, cannot release as %q", profile, entry.Owner, owner)
	}
	return ls.store.Delete(key)
}

// Get returns the current LockEntry for a profile, or an error if none exists.
func (ls *LockStore) Get(profile string) (LockEntry, error) {
	p, err := ls.store.Load(ls.lockKey(profile))
	if err != nil {
		return LockEntry{}, err
	}
	return toLockEntry(p), nil
}

func toLockEntry(p Profile) LockEntry {
	lockedAt, _ := time.Parse(time.RFC3339, p.Vars["LOCKED_AT"])
	expiresAt, _ := time.Parse(time.RFC3339, p.Vars["EXPIRES_AT"])
	return LockEntry{
		Profile:   p.Name,
		Owner:     p.Vars["OWNER"],
		LockedAt:  lockedAt,
		ExpiresAt: expiresAt,
	}
}
