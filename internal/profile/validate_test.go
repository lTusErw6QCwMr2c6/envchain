package profile

import (
	"testing"
)

func TestValidateName_Valid(t *testing.T) {
	valid := []string{"myprofile", "my-profile", "my_profile", "profile.v2", "A1"}
	for _, name := range valid {
		if err := ValidateName(name); err != nil {
			t.Errorf("expected %q to be valid, got: %v", name, err)
		}
	}
}

func TestValidateName_Invalid(t *testing.T) {
	invalid := []string{"", "my profile", "my/profile", "prof@ile", "!bad"}
	for _, name := range invalid {
		if err := ValidateName(name); err == nil {
			t.Errorf("expected %q to be invalid, but got no error", name)
		}
	}
}

func TestValidateVarName_Valid(t *testing.T) {
	valid := []string{"FOO", "_BAR", "foo_bar", "A1", "MY_VAR_123"}
	for _, v := range valid {
		if err := ValidateVarName(v); err != nil {
			t.Errorf("expected %q to be valid, got: %v", v, err)
		}
	}
}

func TestValidateVarName_Invalid(t *testing.T) {
	invalid := []string{"", "1FOO", "FOO BAR", "FOO-BAR", "FOO.BAR"}
	for _, v := range invalid {
		if err := ValidateVarName(v); err == nil {
			t.Errorf("expected %q to be invalid, but got no error", v)
		}
	}
}

func TestValidateProfile_Valid(t *testing.T) {
	p := &Profile{
		Name: "test-profile",
		Vars: []Var{{Name: "FOO"}, {Name: "BAR"}},
	}
	if err := ValidateProfile(p); err != nil {
		t.Fatalf("expected valid profile, got: %v", err)
	}
}

func TestValidateProfile_DuplicateVar(t *testing.T) {
	p := &Profile{
		Name: "dup-test",
		Vars: []Var{{Name: "FOO"}, {Name: "FOO"}},
	}
	err := ValidateProfile(p)
	if err == nil {
		t.Fatal("expected error for duplicate variable")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Errors) != 1 {
		t.Errorf("expected 1 error, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func TestValidateProfile_MultipleErrors(t *testing.T) {
	p := &Profile{
		Name: "bad name!",
		Vars: []Var{{Name: "1INVALID"}, {Name: "OK"}, {Name: "OK"}},
	}
	err := ValidateProfile(p)
	if err == nil {
		t.Fatal("expected errors")
	}
	ve := err.(*ValidationError)
	// Expect: bad name, invalid var name, duplicate OK
	if len(ve.Errors) < 3 {
		t.Errorf("expected at least 3 errors, got %d: %v", len(ve.Errors), ve.Errors)
	}
}

func TestValidationError_Error(t *testing.T) {
	ve := &ValidationError{Errors: []string{"err1", "err2"}}
	msg := ve.Error()
	if msg == "" {
		t.Error("expected non-empty error message")
	}
}
