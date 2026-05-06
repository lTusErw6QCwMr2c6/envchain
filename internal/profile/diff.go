package profile

import "fmt"

// DiffKind represents the type of change in a diff entry.
type DiffKind string

const (
	DiffAdded   DiffKind = "added"
	DiffRemoved DiffKind = "removed"
	DiffChanged DiffKind = "changed"
)

// DiffEntry represents a single variable-level change between two profiles.
type DiffEntry struct {
	Key  string
	Kind DiffKind
	Old  string // empty for added entries
	New  string // empty for removed entries
}

func (d DiffEntry) String() string {
	switch d.Kind {
	case DiffAdded:
		return fmt.Sprintf("+ %s=%s", d.Key, d.New)
	case DiffRemoved:
		return fmt.Sprintf("- %s=%s", d.Key, d.Old)
	case DiffChanged:
		return fmt.Sprintf("~ %s: %s -> %s", d.Key, d.Old, d.New)
	default:
		return fmt.Sprintf("? %s", d.Key)
	}
}

// DiffProfiles compares two profiles and returns a list of variable-level
// differences. The comparison is based on variable names and values only;
// metadata fields such as Name, Description, and Parents are ignored.
func DiffProfiles(base, target *Profile) []DiffEntry {
	baseMap := base.ToEnvMap()
	targetMap := target.ToEnvMap()

	var entries []DiffEntry

	// Find removed or changed keys.
	for k, oldVal := range baseMap {
		newVal, ok := targetMap[k]
		if !ok {
			entries = append(entries, DiffEntry{Key: k, Kind: DiffRemoved, Old: oldVal})
		} else if oldVal != newVal {
			entries = append(entries, DiffEntry{Key: k, Kind: DiffChanged, Old: oldVal, New: newVal})
		}
	}

	// Find added keys.
	for k, newVal := range targetMap {
		if _, ok := baseMap[k]; !ok {
			entries = append(entries, DiffEntry{Key: k, Kind: DiffAdded, New: newVal})
		}
	}

	return entries
}
