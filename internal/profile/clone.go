package profile

import "fmt"

// CloneOptions controls how a profile is cloned to a new name.
type CloneOptions struct {
	// Overwrite allows the destination profile to be replaced if it already exists.
	Overwrite bool
	// StripParents removes the parent chain from the cloned profile.
	StripParents bool
}

// CloneProfile creates a deep copy of srcName into dstName within the same store.
// Unlike CopyProfile, CloneProfile preserves the full parent chain by default and
// allows selectively stripping it via CloneOptions.StripParents.
func CloneProfile(st *Store, srcName, dstName string, opts CloneOptions) error {
	if err := ValidateName(dstName); err != nil {
		return fmt.Errorf("clone: invalid destination name %q: %w", dstName, err)
	}

	src, err := st.Load(srcName)
	if err != nil {
		return fmt.Errorf("clone: load source %q: %w", srcName, err)
	}

	if !opts.Overwrite {
		_, err := st.Load(dstName)
		if err == nil {
			return fmt.Errorf("clone: destination profile %q already exists (use --overwrite to replace)", dstName)
		}
	}

	// Deep-copy vars
	vars := make([]Var, len(src.Vars))
	copy(vars, src.Vars)

	// Copy parents unless stripped
	parents := src.Parents
	if opts.StripParents {
		parents = nil
	}

	dst := &Profile{
		Name:    dstName,
		Vars:    vars,
		Parents: parents,
	}

	if err := st.Save(dst); err != nil {
		return fmt.Errorf("clone: save destination %q: %w", dstName, err)
	}
	return nil
}
