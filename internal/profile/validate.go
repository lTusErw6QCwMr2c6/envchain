package profile

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// validNameRe matches profile names: alphanumeric, dashes, underscores, dots.
	validNameRe = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	// validVarRe matches environment variable names.
	validVarRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

// ValidationError holds one or more validation failures.
type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.Join(e.Errors, "; "))
}

// ValidateName checks that a profile name is well-formed.
func ValidateName(name string) error {
	if name == "" {
		return fmt.Errorf("profile name must not be empty")
	}
	if !validNameRe.MatchString(name) {
		return fmt.Errorf("profile name %q contains invalid characters (allowed: a-z, A-Z, 0-9, _, -, .)", name)
	}
	return nil
}

// ValidateVarName checks that an environment variable name is well-formed.
func ValidateVarName(v string) error {
	if v == "" {
		return fmt.Errorf("variable name must not be empty")
	}
	if !validVarRe.MatchString(v) {
		return fmt.Errorf("variable name %q is not a valid identifier", v)
	}
	return nil
}

// ValidateProfile runs all validation rules against p and returns a
// *ValidationError if any rule is violated, or nil on success.
func ValidateProfile(p *Profile) error {
	var errs []string

	if err := ValidateName(p.Name); err != nil {
		errs = append(errs, err.Error())
	}

	seen := make(map[string]struct{}, len(p.Vars))
	for _, v := range p.Vars {
		if err := ValidateVarName(v.Name); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if _, dup := seen[v.Name]; dup {
			errs = append(errs, fmt.Sprintf("duplicate variable %q", v.Name))
		}
		seen[v.Name] = struct{}{}
	}

	if len(errs) > 0 {
		return &ValidationError{Errors: errs}
	}
	return nil
}
