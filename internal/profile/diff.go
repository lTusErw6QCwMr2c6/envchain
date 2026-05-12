package profile

import "sort"

// DiffResult describes the difference between two variable maps.
type DiffResult struct {
	Added   []string // keys present in new but not old
	Removed []string // keys present in old but not new
	Changed []string // keys present in both but with different values
}

// DiffProfiles computes the diff between two profiles' variable sets.
func DiffProfiles(old, new *Profile) DiffResult {
	oldMap := make(map[string]string, len(old.Vars))
	for _, v := range old.Vars {
		oldMap[v.Key] = v.Value
	}
	newMap := make(map[string]string, len(new.Vars))
	for _, v := range new.Vars {
		newMap[v.Key] = v.Value
	}
	return diffMaps(oldMap, newMap)
}

// diffMaps is the shared implementation used by DiffProfiles and DiffPin.
func diffMaps(oldMap, newMap map[string]string) DiffResult {
	var result DiffResult

	for k, newVal := range newMap {
		if oldVal, ok := oldMap[k]; !ok {
			result.Added = append(result.Added, k)
		} else if oldVal != newVal {
			result.Changed = append(result.Changed, k)
		}
	}
	for k := range oldMap {
		if _, ok := newMap[k]; !ok {
			result.Removed = append(result.Removed, k)
		}
	}

	sort.Strings(result.Added)
	sort.Strings(result.Removed)
	sort.Strings(result.Changed)
	return result
}
