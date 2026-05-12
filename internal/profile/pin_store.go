package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const pinDir = "pins"

// PinStore persists and retrieves pinned profiles.
type PinStore struct {
	dir string
}

// NewPinStore returns a PinStore rooted at baseDir/pins.
func NewPinStore(baseDir string) (*PinStore, error) {
	dir := filepath.Join(baseDir, pinDir)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("pin store: mkdir: %w", err)
	}
	return &PinStore{dir: dir}, nil
}

func (ps *PinStore) pinPath(name, label string) string {
	return filepath.Join(ps.dir, name+"."+label+".json")
}

// Save persists a pin under the given label.
func (ps *PinStore) Save(label string, pin *PinnedProfile) error {
	data, err := json.MarshalIndent(pin, "", "  ")
	if err != nil {
		return fmt.Errorf("pin store: marshal: %w", err)
	}
	path := ps.pinPath(pin.Name, label)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("pin store: write %s: %w", path, err)
	}
	return nil
}

// Load retrieves a previously saved pin by profile name and label.
func (ps *PinStore) Load(name, label string) (*PinnedProfile, error) {
	path := ps.pinPath(name, label)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("pin store: pin %q/%q not found", name, label)
		}
		return nil, fmt.Errorf("pin store: read %s: %w", path, err)
	}
	var pin PinnedProfile
	if err := json.Unmarshal(data, &pin); err != nil {
		return nil, fmt.Errorf("pin store: unmarshal: %w", err)
	}
	return &pin, nil
}

// Delete removes a saved pin.
func (ps *PinStore) Delete(name, label string) error {
	path := ps.pinPath(name, label)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("pin store: delete %s: %w", path, err)
	}
	return nil
}
