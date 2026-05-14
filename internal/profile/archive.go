package profile

import (
	"fmt"
	"time"
)

// ArchivedProfile holds a profile along with archival metadata.
type ArchivedProfile struct {
	Profile   Profile
	ArchivedAt time.Time
	ArchivedBy string
	Reason     string
}

// ArchiveStore persists archived profiles separately from active ones.
type ArchiveStore interface {
	SaveArchive(name string, a ArchivedProfile) error
	LoadArchive(name string) (ArchivedProfile, error)
	ListArchives() ([]string, error)
	DeleteArchive(name string) error
}

// ArchiveProfile moves a profile from the active store into the archive store.
// The profile is removed from the active store on success.
func ArchiveProfile(src Store, dst ArchiveStore, name, archivedBy, reason string) error {
	p, err := src.Load(name)
	if err != nil {
		return fmt.Errorf("archive: load profile %q: %w", name, err)
	}

	a := ArchivedProfile{
		Profile:    p,
		ArchivedAt: time.Now().UTC(),
		ArchivedBy: archivedBy,
		Reason:     reason,
	}

	if err := dst.SaveArchive(name, a); err != nil {
		return fmt.Errorf("archive: save archive %q: %w", name, err)
	}

	if err := src.Delete(name); err != nil {
		return fmt.Errorf("archive: delete active profile %q: %w", name, err)
	}

	return nil
}

// RestoreProfile moves a profile from the archive store back into the active store.
// Returns an error if the destination profile already exists and overwrite is false.
func RestoreProfile(src ArchiveStore, dst Store, name string, overwrite bool) (ArchivedProfile, error) {
	a, err := src.LoadArchive(name)
	if err != nil {
		return ArchivedProfile{}, fmt.Errorf("restore: load archive %q: %w", name, err)
	}

	if !overwrite {
		if _, err := dst.Load(name); err == nil {
			return ArchivedProfile{}, fmt.Errorf("restore: profile %q already exists in active store", name)
		}
	}

	if err := dst.Save(a.Profile); err != nil {
		return ArchivedProfile{}, fmt.Errorf("restore: save profile %q: %w", name, err)
	}

	if err := src.DeleteArchive(name); err != nil {
		return ArchivedProfile{}, fmt.Errorf("restore: delete archive %q: %w", name, err)
	}

	return a, nil
}
