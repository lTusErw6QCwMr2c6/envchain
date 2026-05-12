package profile

import (
	"fmt"
	"time"
)

// Snapshot represents a point-in-time capture of a profile's variables.
type Snapshot struct {
	ProfileName string            `json:"profile_name"`
	CapturedAt  time.Time         `json:"captured_at"`
	Vars        map[string]string `json:"vars"`
}

// SnapshotProfile creates a Snapshot from the given profile using the provided
// env map (typically resolved secrets). The snapshot is a deep copy and does
// not hold references to the original profile.
func SnapshotProfile(p Profile, envMap map[string]string) Snapshot {
	vars := make(map[string]string, len(envMap))
	for k, v := range envMap {
		vars[k] = v
	}
	return Snapshot{
		ProfileName: p.Name,
		CapturedAt:  time.Now().UTC(),
		Vars:        vars,
	}
}

// DiffSnapshot compares a Snapshot against a current env map and returns the
// keys that were added, removed, or changed relative to the snapshot.
func DiffSnapshot(snap Snapshot, current map[string]string) (added, removed, changed []string) {
	for k, v := range current {
		if old, ok := snap.Vars[k]; !ok {
			added = append(added, k)
		} else if old != v {
			changed = append(changed, k)
		}
	}
	for k := range snap.Vars {
		if _, ok := current[k]; !ok {
			removed = append(removed, k)
		}
	}
	return
}

// String returns a human-readable summary of the snapshot.
func (s Snapshot) String() string {
	return fmt.Sprintf("Snapshot(%s @ %s, %d vars)",
		s.ProfileName, s.CapturedAt.Format(time.RFC3339), len(s.Vars))
}
