package profile

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const defaultDir = ".envchain"
const fileExt = ".toml"

// Store persists and retrieves profiles from a directory on disk.
type Store struct {
	Dir string
}

// NewStore returns a Store rooted at dir, creating the directory if needed.
func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("creating profile store dir: %w", err)
	}
	return &Store{Dir: dir}, nil
}

// DefaultStore returns a Store under $HOME/.envchain.
func DefaultStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return NewStore(filepath.Join(home, defaultDir))
}

// Save writes the profile to disk, overwriting any existing file.
func (s *Store) Save(p *Profile) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalid profile %q: %w", p.Name, err)
	}
	path := s.path(p.Name)
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("opening profile file: %w", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(p)
}

// Load reads a profile by name from disk.
func (s *Store) Load(name string) (*Profile, error) {
	var p Profile
	path := s.path(name)
	if _, err := toml.DecodeFile(path, &p); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("profile %q not found", name)
		}
		return nil, fmt.Errorf("decoding profile %q: %w", name, err)
	}
	return &p, nil
}

// List returns the names of all stored profiles.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.Dir)
	if err != nil {
		return nil, fmt.Errorf("reading profile store: %w", err)
	}
	var names []string
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == fileExt {
			names = append(names, e.Name()[:len(e.Name())-len(fileExt)])
		}
	}
	return names, nil
}

func (s *Store) path(name string) string {
	return filepath.Join(s.Dir, name+fileExt)
}
