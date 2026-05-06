package profile

import (
	"fmt"
	"strings"
)

// LintWarning represents a non-fatal advisory about a profile.
type LintWarning struct {
	Field   string
	Message string
}

func (w LintWarning) String() string {
	return fmt.Sprintf("[%s] %s", w.Field, w.Message)
}

// LintResult holds the outcome of linting a profile.
type LintResult struct {
	Profile  string
	Warnings []LintWarning
}

// OK returns true when there are no warnings.
func (r *LintResult) OK() bool { return len(r.Warnings) == 0 }

// LintProfile inspects p for common issues and returns a LintResult.
// Unlike ValidateProfile, lint rules are advisory and do not block saving.
func LintProfile(p *Profile) *LintResult {
	res := &LintResult{Profile: p.Name}

	if len(p.Vars) == 0 {
		res.Warnings = append(res.Warnings, LintWarning{
			Field:   "vars",
			Message: "profile has no variables defined",
		})
	}

	for _, v := range p.Vars {
		if v.Name != strings.ToUpper(v.Name) {
			res.Warnings = append(res.Warnings, LintWarning{
				Field:   "vars." + v.Name,
				Message: fmt.Sprintf("variable %q is not uppercase (conventional style)", v.Name),
			})
		}
		if strings.HasPrefix(v.Name, "_") {
			res.Warnings = append(res.Warnings, LintWarning{
				Field:   "vars." + v.Name,
				Message: fmt.Sprintf("variable %q starts with underscore (reserved by convention)", v.Name),
			})
		}
	}

	if len(p.Parents) > 5 {
		res.Warnings = append(res.Warnings, LintWarning{
			Field:   "parents",
			Message: fmt.Sprintf("profile chains %d parents; consider flattening for clarity", len(p.Parents)),
		})
	}

	return res
}
