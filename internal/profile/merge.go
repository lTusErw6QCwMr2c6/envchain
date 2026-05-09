package profile

// MergeResult holds the outcome of merging two profiles.
type MergeResult struct {
	// Profile is the resulting merged profile.
	Profile *Profile
	// Overwritten contains variable names whose values were overwritten by src.
	Overwritten []string
	// Added contains variable names that were added from src.
	Added []string
}

// MergeProfiles merges src into dst, returning a new Profile and a MergeResult
// describing what changed. Variables in src take precedence over dst when
// overwrite is true; otherwise existing variables in dst are preserved.
//
// The returned Profile is a shallow copy — the original profiles are not
// mutated.
func MergeProfiles(dst, src *Profile, overwrite bool) (*Profile, MergeResult) {
	result := MergeResult{}

	// Start from a copy of dst vars.
	mergedVars := make([]EnvVar, len(dst.Vars))
	copy(mergedVars, dst.Vars)

	// Build an index of existing variable names for quick lookup.
	index := make(map[string]int, len(mergedVars))
	for i, v := range mergedVars {
		index[v.Name] = i
	}

	for _, sv := range src.Vars {
		if i, exists := index[sv.Name]; exists {
			if overwrite {
				mergedVars[i] = sv
				result.Overwritten = append(result.Overwritten, sv.Name)
			}
			// If not overwriting, keep dst value — no action needed.
		} else {
			mergedVars = append(mergedVars, sv)
			index[sv.Name] = len(mergedVars) - 1
			result.Added = append(result.Added, sv.Name)
		}
	}

	merged := &Profile{
		Name:    dst.Name,
		Parents: dst.Parents,
		Vars:    mergedVars,
	}
	result.Profile = merged
	return merged, result
}
