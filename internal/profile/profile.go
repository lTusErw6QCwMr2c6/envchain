package profile

import (
	"errors"
	"fmt"
	"regexp"
)

var validName = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// Var represents a single environment variable entry in a profile.
type Var struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Profile holds a named set of environment variables and an optional chain
// of parent profile names whose variables are merged before this profile's own.
type Profile struct {
	Name  string   `json:"name"`
	Chain []string `json:"chain,omitempty"`
	Vars  []Var    `json:"vars"`
}

// Validate returns an error if the profile has an invalid name or duplicate keys.
func (p *Profile) Validate() error {
	if !validName.MatchString(p.Name) {
		return fmt.Errorf("invalid profile name %q: must match %s", p.Name, validName)
	}
	seen := map[string]bool{}
	for _, v := range p.Vars {
		if seen[v.Key] {
			return fmt.Errorf("duplicate variable key %q in profile %q", v.Key, p.Name)
		}
		seen[v.Key] = true
	}
	return nil
}

// ToEnvMap converts the profile's variables to a plain map.
func (p *Profile) ToEnvMap() (map[string]string, error) {
	if err := p.Validate(); err != nil {
		return nil, err
	}
	m := make(map[string]string, len(p.Vars))
	for _, v := range p.Vars {
		m[v.Key] = v.Value
	}
	return m, nil
}

// ErrProfileNotFound is returned when a requested profile does not exist.
var ErrProfileNotFound = errors.New("profile not found")
