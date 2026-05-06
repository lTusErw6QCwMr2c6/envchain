package profile

import "fmt"

// CopyOptions configures the behaviour of CopyProfile.
type CopyOptions struct {
	// Overwrite allows the destination profile to be replaced if it already exists.
	Overwrite bool
	// ResetParents clears the parent chain on the copied profile.
	ResetParents bool
}

// CopyProfile duplicates src into a new profile named dst inside the given
// Store. Variable values are shallow-copied; the original profile is not
// modified.
func CopyProfile(st *Store, src, dst string, opts CopyOptions) (*Profile, error) {
	if err := ValidateName(dst); err != nil {
		return nil, fmt.Errorf("invalid destination name %q: %w", dst, err)
	}

	original, err := st.Load(src)
	if err != nil {
		return nil, fmt.Errorf("load source profile %q: %w", src, err)
	}

	if !opts.Overwrite {
		if _, err := st.Load(dst); err == nil {
			return nil, fmt.Errorf("destination profile %q already exists; use Overwrite to replace it", dst)
		}
	}

	copy := &Profile{
		Name: dst,
		Vars: make([]Var, len(original.Vars)),
		Parents: make([]string, len(original.Parents)),
	}

	for i, v := range original.Vars {
		copy.Vars[i] = Var{Name: v.Name, Value: v.Value}
	}

	if !opts.ResetParents {
		for i, p := range original.Parents {
			copy.Parents[i] = p
		}
	}

	if err := st.Save(copy); err != nil {
		return nil, fmt.Errorf("save destination profile %q: %w", dst, err)
	}

	return copy, nil
}

// RenameProfile copies src to dst and then deletes the original.
func RenameProfile(st *Store, src, dst string, opts CopyOptions) (*Profile, error) {
	cp, err := CopyProfile(st, src, dst, opts)
	if err != nil {
		return nil, fmt.Errorf("rename (copy phase): %w", err)
	}

	if err := st.Delete(src); err != nil {
		return nil, fmt.Errorf("rename (delete source %q): %w", src, err)
	}

	return cp, nil
}
