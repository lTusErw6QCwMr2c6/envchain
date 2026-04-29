// Package profile manages environment variable profiles for envchain.
package profile

import (
	"errors"
	"regexp"
)

var validNameRe = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

// ErrInvalidName is returned when a profile name contains invalid characters.
var ErrInvalidName = errors.New("profile name must match [a-zA-Z0-9_-]")

// ErrDuplicateVar is returned when a variable is defined more than once.
var ErrDuplicateVar = errors.New("duplicate variable in profile")

// Var represents a single environment variable entry in a profile.
type Var struct {
	Key    string `toml:"key"`
	Value  string `toml:"value,omitempty"`
	Secret bool   `toml:"secret,omitempty"`
	Ref    string `toml:"ref,omitempty"` // secret store reference e.g. "vault:secret/myapp#DB_PASS"
}

// Profile holds a named collection of environment variables.
type Profile struct {
	Name    string   `toml:"name"`
	Extends []string `toml:"extends,omitempty"`
	Vars    []Var    `toml:"vars"`
}

// Validate checks that the profile is well-formed.
func (p *Profile) Validate() error {
	if !validNameRe.MatchString(p.Name) {
		return ErrInvalidName
	}
	seen := make(map[string]struct{}, len(p.Vars))
	for _, v := range p.Vars {
		if _, exists := seen[v.Key]; exists {
			return ErrDuplicateVar
		}
		seen[v.Key] = struct{}{}
	}
	return nil
}

// ToEnvMap converts the profile's vars into a key→value map.
// Secret vars with no resolved Value are included with an empty string.
func (p *Profile) ToEnvMap() map[string]string {
	env := make(map[string]string, len(p.Vars))
	for _, v := range p.Vars {
		env[v.Key] = v.Value
	}
	return env
}
