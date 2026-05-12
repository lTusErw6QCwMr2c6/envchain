package profile

import "fmt"

// RenameOptions controls behaviour when renaming a profile.
type RenameOptions struct {
	// Overwrite allows the destination profile to be replaced if it already exists.
	Overwrite bool
}

// RenameProfile renames an existing profile from oldName to newName.
// It loads the profile, updates its name, saves it under the new name, and
// removes the old entry. If the destination already exists and Overwrite is
// false an error is returned.
func RenameProfile(st *Store, oldName, newName string, opts RenameOptions) error {
	if err := ValidateName(newName); err != nil {
		return fmt.Errorf("rename: invalid destination name %q: %w", newName, err)
	}

	src, err := st.Load(oldName)
	if err != nil {
		return fmt.Errorf("rename: load %q: %w", oldName, err)
	}

	if !opts.Overwrite {
		if _, err := st.Load(newName); err == nil {
			return fmt.Errorf("rename: destination profile %q already exists (use --overwrite to replace)", newName)
		}
	}

	dst := *src
	dst.Name = newName

	if err := st.Save(&dst); err != nil {
		return fmt.Errorf("rename: save %q: %w", newName, err)
	}

	if err := st.Delete(oldName); err != nil {
		// Best-effort rollback: remove the newly saved profile so we don't
		// leave the store in an inconsistent state.
		_ = st.Delete(newName)
		return fmt.Errorf("rename: delete old profile %q: %w", oldName, err)
	}

	return nil
}
