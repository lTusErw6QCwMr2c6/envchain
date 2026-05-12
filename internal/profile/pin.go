package profile

import (
	"fmt"
	"time"
)

// PinnedProfile records a specific version of a profile at a point in time.
type PinnedProfile struct {
	Name      string            `json:"name"`
	PinnedAt  time.Time         `json:"pinned_at"`
	PinnedBy  string            `json:"pinned_by"`
	Vars      map[string]string `json:"vars"`
	Parents   []string          `json:"parents"`
}

// PinProfile captures the current state of a profile as a named pin.
func PinProfile(s *Store, name, pinnedBy string) (*PinnedProfile, error) {
	p, err := s.Load(name)
	if err != nil {
		return nil, fmt.Errorf("pin: load profile %q: %w", name, err)
	}

	vars := make(map[string]string, len(p.Vars))
	for _, v := range p.Vars {
		vars[v.Key] = v.Value
	}

	pin := &PinnedProfile{
		Name:     name,
		PinnedAt: time.Now().UTC(),
		PinnedBy: pinnedBy,
		Vars:     vars,
		Parents:  append([]string(nil), p.Parents...),
	}
	return pin, nil
}

// DiffPin compares a pinned snapshot against the current profile state.
// It returns a DiffResult describing what changed since the pin was taken.
func DiffPin(s *Store, pin *PinnedProfile) (DiffResult, error) {
	p, err := s.Load(pin.Name)
	if err != nil {
		return DiffResult{}, fmt.Errorf("diff pin: load profile %q: %w", pin.Name, err)
	}

	current := make(map[string]string, len(p.Vars))
	for _, v := range p.Vars {
		current[v.Key] = v.Value
	}

	return diffMaps(pin.Vars, current), nil
}
